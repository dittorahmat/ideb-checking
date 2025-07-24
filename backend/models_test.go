package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestTableName(t *testing.T) {
	req := Request{}
	assert.Equal(t, "getDebtorExactIndividual", req.TableName())
}

func TestCorporateRequestTableName(t *testing.T) {
	corpReq := CorporateRequest{}
	assert.Equal(t, "getDebtorExactCorporate", corpReq.TableName())
}

func TestGetIdebTableName(t *testing.T) {
	getIdeb := GetIdeb{}
	assert.Equal(t, "get_idebs", getIdeb.TableName())
}