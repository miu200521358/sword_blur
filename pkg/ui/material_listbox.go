package ui

import (
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/walk/pkg/walk"
)

type MaterialListBox struct {
	*walk.ListBox
}

func NewMaterialListBox(
	parent walk.Container,
) (*MaterialListBox, error) {
	lb, err := walk.NewListBox(parent)
	if err != nil {
		return nil, err
	}
	lb.SetMinMaxSize(walk.Size{Width: -1, Height: 100}, walk.Size{Width: -1, Height: 200})

	return &MaterialListBox{lb}, nil
}

func (lb *MaterialListBox) SetMaterials(
	materials *pmx.Materials, funcChange func(indexes []int),
) {
	model := NewMaterialListModel(materials)
	lb.SetModel(model)
	lb.SetSelectedIndexes(materials.GetIndexes())
	lb.SelectedIndexesChanged().Attach(func() {
		funcChange(lb.SelectedIndexes())
	})
}

type MaterialItem struct {
	name  string
	value int
}

type MaterialListModel struct {
	walk.ListModelBase
	items []*MaterialItem
}

func NewMaterialListModel(materials *pmx.Materials) *MaterialListModel {
	m := &MaterialListModel{items: make([]*MaterialItem, materials.Len())}

	for i := range materials.Len() {
		material := materials.Get(i)

		m.items[i] = &MaterialItem{material.Name, material.Index}
	}

	return m
}

func (m *MaterialListModel) ItemCount() int {
	return len(m.items)
}

func (m *MaterialListModel) Value(index int) string {
	return m.items[index].name
}
