package generator

import (
	"testing"

	mock_passwordstore "unique-pass-gen/pkg/passwordstore/mocks"

	"github.com/golang/mock/gomock"
)

func hasDuplicateChars(s string) bool {
	seen := make(map[rune]bool)
	for _, r := range s {
		if seen[r] {
			return true
		}

		seen[r] = true
	}

	return false
}

func containsFromSet(password string, set []rune) bool {
	for _, r := range password {
		for _, s := range set {
			if r == s {
				return true
			}
		}
	}

	return false
}

func TestUniquePasswordGenerator_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	tests := []struct {
		name           string
		options        Options
		expectErr      bool
		expectedLength int
	}{
		{
			name:           "valid digits+lower+upper",
			options:        NewOptions(WithLength(10), WithDigits(), WithLowerC(), WithUpperC()),
			expectErr:      false,
			expectedLength: 10,
		},
		{
			name:           "valid digits+lower",
			options:        NewOptions(WithLength(5), WithDigits(), WithLowerC()),
			expectErr:      false,
			expectedLength: 5,
		},
		{
			name:           "valid digits+upper",
			options:        NewOptions(WithLength(10), WithDigits(), WithUpperC()),
			expectErr:      false,
			expectedLength: 10,
		},
		{
			name:           "valid lower+upper",
			options:        NewOptions(WithLength(10), WithLowerC(), WithUpperC()),
			expectErr:      false,
			expectedLength: 10,
		},
		{
			name:           "valid digits",
			options:        NewOptions(WithLength(10), WithDigits()),
			expectErr:      false,
			expectedLength: 10,
		},
		{
			name:           "valid lower",
			options:        NewOptions(WithLength(20), WithLowerC()),
			expectErr:      false,
			expectedLength: 20,
		},
		{
			name:           "valid upper",
			options:        NewOptions(WithLength(20), WithUpperC()),
			expectErr:      false,
			expectedLength: 20,
		},
	}

	for _, tt := range tests {
		mockStore := mock_passwordstore.NewMockPasswordStore(ctrl)

		if !tt.expectErr {
			mockStore.EXPECT().Exists(gomock.Any()).Return(false)
			mockStore.EXPECT().Add(gomock.Any())
			mockStore.EXPECT().Get().Return([]string{})
		}

		gen := NewGenerator(mockStore)

		password, err := gen.UniquePasswordGenerator(tt.options)
		if tt.expectErr {
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			return
		}

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(password) != tt.expectedLength {
			t.Errorf("expected length %d, got %d", tt.expectedLength, len(password))
		}
	}
}

func TestUniquePasswordGenerator_NoOptionsSelected(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := mock_passwordstore.NewMockPasswordStore(ctrl)
	gen := NewGenerator(mockStore)
	options := NewOptions(WithLength(5))

	_, err := gen.UniquePasswordGenerator(options)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUniquePasswordGenerator_NoDuplicatesInPassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := mock_passwordstore.NewMockPasswordStore(ctrl)
	mockStore.EXPECT().Exists(gomock.Any()).Return(false)
	mockStore.EXPECT().Add(gomock.Any())
	mockStore.EXPECT().Get().Return([]string{})

	gen := NewGenerator(mockStore)
	options := NewOptions(WithLength(10), WithDigits(), WithLowerC(), WithUpperC())

	password, err := gen.UniquePasswordGenerator(options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hasDuplicateChars(password) {
		t.Errorf("password has duplicate characters: %s", password)
	}
}

func TestUniquePasswordGenerator_ContainsAllSelectedSets(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := mock_passwordstore.NewMockPasswordStore(ctrl)
	mockStore.EXPECT().Exists(gomock.Any()).Return(false)
	mockStore.EXPECT().Add(gomock.Any())
	mockStore.EXPECT().Get().Return([]string{})

	gen := NewGenerator(mockStore)
	options := NewOptions(WithLength(10), WithDigits(), WithLowerC(), WithUpperC())

	password, err := gen.UniquePasswordGenerator(options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsFromSet(password, digits) {
		t.Errorf("password missing digit: %s", password)
	}

	if !containsFromSet(password, lowerC) {
		t.Errorf("password missing lowercase: %s", password)
	}

	if !containsFromSet(password, upperC) {
		t.Errorf("password missing uppercase: %s", password)
	}
}

func TestPasswordsAreUniqueAcrossGenerations(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := mock_passwordstore.NewMockPasswordStore(ctrl)
	mockStore.EXPECT().Exists(gomock.Any()).Return(false).AnyTimes()
	mockStore.EXPECT().Add(gomock.Any()).AnyTimes()
	mockStore.EXPECT().Get().Return([]string{}).AnyTimes()

	gen := NewGenerator(mockStore)
	options := NewOptions(WithLength(10), WithDigits(), WithLowerC(), WithUpperC())

	seen := make(map[string]bool)

	for i := 0; i < 100; i++ {
		password, err := gen.UniquePasswordGenerator(options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if seen[password] {
			t.Fatalf("duplicate password generated: %s", password)
		}

		seen[password] = true
	}
}
