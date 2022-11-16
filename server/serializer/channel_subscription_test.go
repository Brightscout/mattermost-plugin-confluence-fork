package serializer

import (
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestFormattedSubscriptionList(t *testing.T) {
	tests := map[string]struct {
		subscription map[string]Subscription
		result       string
	}{
		"space subscription": {
			subscription: map[string]Subscription{
				"test": &SpaceSubscription{
					SpaceKey: "TS",
					BaseSubscription: BaseSubscription{
						Alias:     "test",
						BaseURL:   "https://test.com",
						ChannelID: "testChannelID",
						Events:    []string{CommentRemovedEvent, CommentUpdatedEvent},
					},
				},
			},
			result: "#### Space Subscriptions \n| Name | Base URL | Space Key | Events|\n| :----|:--------| :--------| :-----|\n|test|https://test.com|TS|Comment Remove, Comment Update|",
		},
		"page subscription": {
			subscription: map[string]Subscription{
				"abc": &PageSubscription{
					PageID: "12345",
					BaseSubscription: BaseSubscription{
						Alias:     "abc",
						BaseURL:   "https://test.com",
						ChannelID: "testChannelID",
						Events:    []string{CommentCreatedEvent, CommentUpdatedEvent},
					},
				},
			},
			result: "#### Page Subscriptions \n| Name | Base URL | Page ID | Events|\n| :----|:--------| :--------| :-----|\n|abc|https://test.com|12345|Comment Create, Comment Update|",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()

			formattedList := FormattedSubscriptionList(tc.subscription)
			assert.Equal(t, tc.result, formattedList)
		})
	}
}

func TestFormattedOldSubscriptionList(t *testing.T) {
	tests := map[string]struct {
		subscription []Subscription
		result       string
	}{
		"space subscription": {
			subscription: []Subscription{
				&SpaceSubscription{
					SpaceKey: "TS",
					BaseSubscription: BaseSubscription{
						Alias:     "test",
						BaseURL:   "https://test.com",
						ChannelID: "testChannelID",
						Events:    []string{CommentRemovedEvent, CommentUpdatedEvent},
					},
				},
			},
			result: "#### Space Subscriptions \n| Name | Base URL | Space Key | Channel ID | Events|\n| :----|:--------| :--------| :-----|\n|test|https://test.com|TS|testChannelID|Comment Remove, Comment Update|",
		},
		"page subscription": {
			subscription: []Subscription{
				&PageSubscription{
					PageID: "12345",
					BaseSubscription: BaseSubscription{
						Alias:     "abc",
						BaseURL:   "https://test.com",
						ChannelID: "testChannelID",
						Events:    []string{CommentCreatedEvent, CommentUpdatedEvent},
					},
				},
			},
			result: "#### Page Subscriptions \n| Name | Base URL | Page ID | Channel ID | Events|\n| :----|:--------| :--------| :-----|\n|abc|https://test.com|12345|testChannelID|Comment Create, Comment Update|",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()

			formattedList := FormattedOldSubscriptionList(tc.subscription)
			assert.Equal(t, tc.result, formattedList)
		})
	}
}

func TestGetOldFormattedSubscription(t *testing.T) {
	tests := map[string]struct {
		subscription Subscription
		result       string
	}{
		"space subscription": {
			subscription: &SpaceSubscription{
				SpaceKey: "TS",
				BaseSubscription: BaseSubscription{
					Alias:     "test",
					BaseURL:   "https://test.com",
					ChannelID: "testChannelID",
					Events:    []string{CommentRemovedEvent, CommentUpdatedEvent},
				},
			},
			result: "\n|test|https://test.com|TS|testChannelID|Comment Remove, Comment Update|",
		},
		"page subscription": {
			subscription: &PageSubscription{
				PageID: "12345",
				BaseSubscription: BaseSubscription{
					Alias:     "abc",
					BaseURL:   "https://test.com",
					ChannelID: "testChannelID",
					Events:    []string{CommentCreatedEvent, CommentUpdatedEvent},
				},
			},
			result: "\n|abc|https://test.com|12345|testChannelID|Comment Create, Comment Update|",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()

			formattedSubscription := tc.subscription.GetOldFormattedSubscription()
			assert.Equal(t, tc.result, formattedSubscription)
		})
	}
}

func TestGetFormattedSubscription(t *testing.T) {
	tests := map[string]struct {
		subscription Subscription
		result       string
	}{
		"space subscription": {
			subscription: &SpaceSubscription{
				SpaceKey: "TS",
				BaseSubscription: BaseSubscription{
					Alias:     "test",
					BaseURL:   "https://test.com",
					ChannelID: "testChannelID",
					Events:    []string{CommentRemovedEvent, CommentUpdatedEvent},
				},
			},
			result: "\n|test|https://test.com|TS|Comment Remove, Comment Update|",
		},
		"page subscription": {
			subscription: &PageSubscription{
				PageID: "12345",
				BaseSubscription: BaseSubscription{
					Alias:     "abc",
					BaseURL:   "https://test.com",
					ChannelID: "testChannelID",
					Events:    []string{CommentCreatedEvent, CommentUpdatedEvent},
				},
			},
			result: "\n|abc|https://test.com|12345|Comment Create, Comment Update|",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()

			formattedSubscription := tc.subscription.GetFormattedSubscription()
			assert.Equal(t, tc.result, formattedSubscription)
		})
	}
}
