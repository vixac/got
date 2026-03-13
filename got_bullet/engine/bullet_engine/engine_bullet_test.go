package bullet_engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"vixac.com/got/engine"
)

func TestCreateBuckWithOverrideSettings(t *testing.T) {
	mock_client := BuildTestClient()
	sut, err := NewEngineBullet(mock_client)
	assert.NoError(t, err)

	flags := []string{"flag1", "flag2"}

	tag1Literal := engine.TagLiteral{
		Display: "tag1",
		Token:   "",
	}
	tag1 := engine.Tag{
		Identifier: nil,
		Literal:    &tag1Literal,
	}
	tags := []engine.Tag{
		tag1,
	}

	var overrideId int32 = 1360
	var longForm string = "This is a long form text entry."
	override := engine.CreateOverrideSettings{
		OverrideId:  &overrideId,
		UpdatedDate: "2026-01-14T18:39:21.429465Z",
		CreatedDate: "2026-01-13T18:39:21.429465Z",
		ScheduleDate: &engine.DateTime{
			Special: "n",
		},
		Tags:     tags,
		Flags:    flags,
		LongForm: &longForm,
	}

	buck1Id := overrideId
	req1 := engine.NewCreateBuckRequest(nil, nil, "buck1", engine.Active, &override)
	id, err := sut.CreateBuck(req1)
	assert.NoError(t, err)
	assert.Equal(t, id.IntValue, overrideId)

	//fetch items below -1, which is buck1, (buck2 is 0 as it was added most recently)
	items, err := sut.FetchItemsBelow(&engine.GidLookup{
		Input: "",
	}, false, []engine.GotState{engine.Active}, false)
	assert.NoError(t, err)

	assert.Equal(t, len(items.Result), 1)

	firstItem := items.Result[0]
	assert.Equal(t, firstItem.GotId.IntValue, buck1Id)

	//VX:TODO we either inject now or we change the way this is displayed. assert.Equal(t, firstItem.Created, "58 days ago")
	assert.Equal(t, firstItem.Updated, "2026-01-14")
	assert.Equal(t, firstItem.Deadline, "---Now---")

	item1Tags := firstItem.SummaryObj.Tags
	item1Flags := firstItem.SummaryObj.Flags
	assert.Equal(t, len(item1Tags), 1)
	assert.Equal(t, len(item1Flags), 2)
	assert.Equal(t, item1Tags[0].Literal.Display, "tag1")
	assert.Equal(t, item1Flags["flag1"], true)
	assert.Equal(t, item1Flags["flag2"], true)

	longformRes, err := sut.LongFormStore.LongFormForMany([]int32{buck1Id})
	assert.Equal(t, len(longformRes), 1)
	assert.Equal(t, longformRes[buck1Id], "This is a long form text entry.")

}

func TestCreateCompleteBuck(t *testing.T) {
	mock_client := BuildTestClient()
	sut, err := NewEngineBullet(mock_client)
	assert.NoError(t, err)

	req1 := engine.NewCreateBuckRequest(nil, nil, "buck1", engine.Active, nil)
	//create buck1

	buck1Id := int32(360)
	id, err := sut.CreateBuck(req1)
	assert.NoError(t, err)
	assert.Equal(t, id.IntValue, buck1Id)

	req2 := engine.NewCreateBuckRequest(&engine.GidLookup{
		Input: "0",
	}, nil, "buck1:buck2", engine.Active, nil)

	buck2Id := int32(361)
	//create buck2 under buck1
	id, err = sut.CreateBuck(req2)
	assert.NoError(t, err)
	assert.Equal(t, id.IntValue, buck2Id)

	//fetch items below -1, which is buck1, (buck2 is 0 as it was added most recently)
	items, err := sut.FetchItemsBelow(&engine.GidLookup{
		Input: "-1",
	}, false, []engine.GotState{engine.Active}, false)
	assert.NoError(t, err)
	assert.Equal(t, len(items.Result), 1)
	assert.Equal(t, items.Result[0].GotId.IntValue, buck2Id)

	//now add a complete item, buck3 under buck1
	buck3Id := int32(362)
	req3 := engine.NewCreateBuckRequest(&engine.GidLookup{
		Input: "0",
	}, nil, "buck1:buck2:buck3", engine.Complete, nil)

	id, err = sut.CreateBuck(req3)
	assert.NoError(t, err)
	assert.Equal(t, id.IntValue, buck3Id)
	items, err = sut.FetchItemsBelow(&engine.GidLookup{
		Input: "-2", //fetches from the top
	}, false, []engine.GotState{engine.Active, engine.Complete}, false)

	assert.NoError(t, err)
	assert.Equal(t, len(items.Result), 2) //buck2 and buck3 are present.

	//now we add a complete item UNDER buck3
	buck4Id := int32(363)
	req4 := engine.NewCreateBuckRequest(&engine.GidLookup{
		Input: "0",
	}, nil, "buck1:buck2:buck3:buck4", engine.Complete, nil)

	id, err = sut.CreateBuck(req4)
	assert.NoError(t, err)
	assert.Equal(t, id.IntValue, buck4Id)

	items, err = sut.FetchItemsBelow(&engine.GidLookup{
		Input: "", //fetches from the top
	}, false, []engine.GotState{engine.Active, engine.Complete}, false)

	assert.NoError(t, err)
	assert.Equal(t, len(items.Result), 4) //buck2 and buck3 are present.

}
