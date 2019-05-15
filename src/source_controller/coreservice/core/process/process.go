/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package process

import (
	"fmt"
	
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type processOperation struct {
	dbProxy dal.RDB
}

// New create a new model manager instance
func New(dbProxy dal.RDB) core.ProcessOperation {
	processOps := &processOperation{dbProxy: dbProxy}
	return processOps
}

func (p *processOperation) validateBizID(ctx core.ContextParams, md metadata.Metadata) (int64, error) {
	// extract biz id from metadata
	bizID, err := metadata.BizIDFromMetadata(md)
	if err != nil {
		blog.Errorf("parse biz id from metadata failed, err: %+v", err)
		return 0, err
	}
	
	// avoid unnecessary db query
	if bizID == 0 {
		return 0, fmt.Errorf("bizID invalid, bizID: %d", bizID)
	}

	// check bizID valid
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	count, err := p.dbProxy.Table(common.BKTableNameBaseApp).Find(cond.ToMapStr()).Count(ctx.Context)
	if nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameObjDes, err.Error(), ctx.ReqID)
		return  0, err
	}
	if count < 1 {
		return 0, fmt.Errorf("business not found, id:%d", bizID)
	}
	
	return bizID, nil
}