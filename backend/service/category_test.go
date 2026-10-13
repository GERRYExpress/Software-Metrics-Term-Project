package service

import (
	"reflect"
	"software-metrics-project/model"
	"testing"
)

func TestCategoryService_GetCategories(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		s       *CategoryService
		want    []model.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.s.GetCategories()
			if (err != nil) != tt.wantErr {
				t.Fatalf("CategoryService.GetCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryService.GetCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryService_Create(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		s       *CategoryService
		args    args
		want    model.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.s.Create(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CategoryService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryService_Update(t *testing.T) {
	t.Parallel()
	type args struct {
		id       int
		name     string
		isHidden bool
	}
	tests := []struct {
		name    string
		s       *CategoryService
		args    args
		want    model.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.s.Update(tt.args.id, tt.args.name, tt.args.isHidden)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CategoryService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryService_Delete(t *testing.T) {
	t.Parallel()
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		s       *CategoryService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.s.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("CategoryService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
