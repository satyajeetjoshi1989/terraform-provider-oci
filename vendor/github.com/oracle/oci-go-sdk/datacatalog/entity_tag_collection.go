// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Data Catalog API
//
// Use the Data Catalog APIs to collect, organize, find, access, understand, enrich, and activate technical, business, and operational metadata.
//

package datacatalog

import (
	"github.com/oracle/oci-go-sdk/common"
)

// EntityTagCollection Results of an entity tags listing. Entity tags allow assciation of business terms with entities.
type EntityTagCollection struct {

	// Collection of entity tags.
	Items []EntityTagSummary `mandatory:"true" json:"items"`
}

func (m EntityTagCollection) String() string {
	return common.PointerString(m)
}
