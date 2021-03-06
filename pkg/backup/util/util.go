// Copyright 2016 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetLatestBackupName(names []string) string {
	bnames := FilterAndSortBackups(names)
	if len(bnames) == 0 {
		return ""
	}
	return bnames[len(bnames)-1]
}

func IsBackup(name string) bool {
	return strings.HasSuffix(name, BackupFilenameSuffix)
}

func MakeBackupName(ver string, rev int64) string {
	return fmt.Sprintf("%s_%016x_%s", ver, rev, BackupFilenameSuffix)
}

func GetRev(name string) (int64, error) {
	parts := strings.SplitN(name, "_", 3)
	if len(parts) != 3 {
		return 0, fmt.Errorf("bad backup name: %s", name)
	}

	return strconv.ParseInt(parts[1], 16, 64)
}

func FilterAndSortBackups(names []string) []string {
	bnames := make(backupNames, 0)
	for _, n := range names {
		if !IsBackup(n) {
			continue
		}
		_, err := GetRev(n)
		if err != nil {
			logrus.Errorf("fail to get rev from backup (%s): %v", n, err)
			continue
		}
		bnames = append(bnames, n)
	}

	sort.Sort(bnames)
	return []string(bnames)
}

type backupNames []string

func (bn backupNames) Len() int { return len(bn) }

func (bn backupNames) Less(i, j int) bool {
	ri, err := GetRev(bn[i])
	if err != nil {
		panic(err)
	}
	rj, err := GetRev(bn[j])
	if err != nil {
		panic(err)
	}

	return ri < rj
}

func (bn backupNames) Swap(i, j int) {
	bn[i], bn[j] = bn[j], bn[i]
}

func ToMB(s int64) float64 {
	n := float64(s) / (1024 * 1024)
	// truncate to KB
	sn := fmt.Sprintf("%.3f", n)

	var err error
	n, err = strconv.ParseFloat(sn, 64)
	if err != nil {
		panic("unexpected parse float error")
	}
	return n
}
