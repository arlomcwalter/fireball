package ast

import (
	"fireball/core/types"
)

// Struct

func (s *Struct) Size() int {
	size := 0

	for _, field := range s.Fields {
		size += field.Type.Size()
	}

	return size
}

func (s *Struct) WithoutRange() types.Type {
	return &Struct{
		parent: s.parent,
		Name:   s.Name,
		Fields: s.Fields,
		Type:   s.Type,
	}
}

func (s *Struct) Equals(other types.Type) bool {
	if v, ok := other.(*Struct); ok {
		if len(s.Fields) != len(v.Fields) {
			return false
		}

		for i, field := range s.Fields {
			if !field.Type.Equals(v.Fields[i].Type) {
				return false
			}
		}

		return true
	}

	return false
}

func (s *Struct) CanAssignTo(other types.Type) bool {
	return s.Equals(other)
}

// Enum

func (e *Enum) Size() int {
	return e.Type.Size()
}

func (e *Enum) WithoutRange() types.Type {
	return &Enum{
		parent:    e.parent,
		Name:      e.Name,
		Type:      e.Type,
		InferType: e.InferType,
		Cases:     e.Cases,
	}
}

func (e *Enum) Equals(other types.Type) bool {
	if v, ok := other.(*Enum); ok {
		if !e.Type.Equals(v.Type) || len(e.Cases) != len(v.Cases) {
			return false
		}

		for i, eCase := range e.Cases {
			vCase := v.Cases[i]

			if eCase.Name != vCase.Name || eCase.Value != vCase.Value {
				return false
			}
		}

		return true
	}

	return false
}

func (e *Enum) CanAssignTo(other types.Type) bool {
	return e.Equals(other)
}

// Function

func (f *Func) Size() int {
	return 4
}

func (f *Func) WithoutRange() types.Type {
	return &Func{
		parent:   f.parent,
		Extern:   f.Extern,
		Name:     f.Name,
		Params:   f.Params,
		Variadic: f.Variadic,
		Returns:  f.Returns,
		Body:     f.Body,
	}
}

func (f *Func) Equals(other types.Type) bool {
	if v, ok := other.(*Func); ok {
		if !f.Returns.Equals(v.Returns) {
			return false
		}

		if len(f.Params) != len(v.Params) {
			return false
		}

		for i, param := range f.Params {
			if !param.Type.Equals(v.Params[i].Type) {
				return false
			}
		}

		return true
	}

	return false
}

func (f *Func) CanAssignTo(other types.Type) bool {
	if v, ok := other.(*Func); ok {
		if !f.Returns.CanAssignTo(v.Returns) {
			return false
		}

		if len(f.Params) != len(v.Params) {
			return false
		}

		for i, param := range f.Params {
			if !param.Type.CanAssignTo(v.Params[i].Type) {
				return false
			}
		}

		return true
	}

	return false
}

func (f *Func) String() string {
	return f.Signature(false)
}
