package facebook

// type GetAdvertsRequest struct {
// 	BuyLocation struct {
// 		Latitude  float64 `json:"latitude"`
// 		Longitude float64 `json:"longitude"`
// 	} `json:"buyLocation"`
// 	Count  int `json:"count"`
// 	Params struct {
// 		Bqf struct {
// 			Callsite string `json:"callsite"`
// 			Query    string `json:"query"`
// 		} `json:"bqf"`
// 		BrowseRequestParams struct {
// 			CommerceEnableLocalPickup bool    `json:"commerce_enable_local_pickup"`
// 			CommerceEnableShipping    bool    `json:"commerce_enable_shipping"`
// 			FilterLocationLatitude    float64 `json:"filter_location_latitude"`
// 			FilterLocationLongitude   float64 `json:"filter_location_longitude"`
// 			FilterPriceLowerBound     int     `json:"filter_price_lower_bound"`
// 			FilterPriceUpperBound     int64   `json:"filter_price_upper_bound"`
// 			FilterRadiusKm            int     `json:"filter_radius_km"`
// 		} `json:"browse_request_params"`
// 		CustomRequestParams struct {
// 			SearchVertical string `json:"search_vertical"`
// 			SeoURL         string `json:"seo_url"`
// 			Surface        string `json:"surface"`
// 		} `json:"custom_request_params"`
// 	} `json:"params"`
// 	SavedSearchID    interface{} `json:"savedSearchID"`
// 	SavedSearchQuery string      `json:"savedSearchQuery"`
// 	Scale            int         `json:"scale"`
// 	TopicPageParams  struct {
// 		LocationID string `json:"location_id"`
// 		URL        string `json:"url"`
// 	} `json:"topicPageParams"`
// 	VehicleParams string `json:"vehicleParams"`
// }

type GetAdvertsRequest struct {
	BuyLocation     BuyLocation      `json:"buyLocation"`
	CategoryIDArray []int64          `json:"categoryIDArray"`
	ContextualData  []ContextualData `json:"contextual_data"`
	Count           int              `json:"count"`
	// Cursor                   interface{}      `json:"cursor"`
	FilterSortingParams      *FilterSortingParams `json:"filterSortingParams,omitempty"`
	MarketplaceBrowseContext string               `json:"marketplaceBrowseContext"`
	// NumericVerticalFields        []interface{}          `json:"numericVerticalFields"`
	// NumericVerticalFieldsBetween []interface{}          `json:"numericVerticalFieldsBetween"`
	PriceRange           []int                  `json:"priceRange"`
	Radius               int                    `json:"radius"`
	Scale                int                    `json:"scale"`
	StringVerticalFields []StringVerticalFields `json:"stringVerticalFields"`
	TopicPageParams      TopicPageParams        `json:"topicPageParams"`
}

type BuyLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ContextualData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type StringVerticalFields struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TopicPageParams struct {
	// LocationID string `json:"location_id"`
	URL string `json:"url"`
}

type FilterSortingParams struct {
	SortByFilter string `json:"sort_by_filter,omitempty"`
	SortOrder    string `json:"sort_order,omitempty"`
}

type GetAdvertsResponse struct {
	Data struct {
		Viewer struct {
			MarketplaceFeedStories struct {
				Edges []struct {
					Node struct {
						Typename  string `json:"__typename"`
						StoryType string `json:"story_type"`
						StoryKey  string `json:"story_key"`
						Tracking  string `json:"tracking"`
						Listing   struct {
							Typename            string `json:"__typename"`
							ID                  string `json:"id"`
							PrimaryListingPhoto struct {
								Typename string `json:"__typename"`
								Image    struct {
									URI string `json:"uri"`
								} `json:"image"`
								ID string `json:"id"`
							} `json:"primary_listing_photo"`
							IsMarketplaceListingRenderable string `json:"__isMarketplaceListingRenderable"`
							CreationTime                   int    `json:"creation_time"`
							FormattedPrice                 struct {
								Text string `json:"text"`
							} `json:"formatted_price"`
							Location struct {
								ReverseGeocode struct {
									City     string `json:"city"`
									State    string `json:"state"`
									CityPage struct {
										DisplayName string `json:"display_name"`
										ID          string `json:"id"`
									} `json:"city_page"`
								} `json:"reverse_geocode"`
							} `json:"location"`
							IsMarketplaceListingWithEnhancedUpdateTime string      `json:"__isMarketplaceListingWithEnhancedUpdateTime"`
							CreationTimeOverride                       int         `json:"creation_time_override"`
							IsHidden                                   bool        `json:"is_hidden"`
							IsLive                                     bool        `json:"is_live"`
							IsPending                                  bool        `json:"is_pending"`
							IsSold                                     bool        `json:"is_sold"`
							IsViewerSeller                             bool        `json:"is_viewer_seller"`
							MinListingPrice                            interface{} `json:"min_listing_price"`
							MaxListingPrice                            interface{} `json:"max_listing_price"`
							MarketplaceListingCategoryID               string      `json:"marketplace_listing_category_id"`
							MarketplaceListingTitle                    string      `json:"marketplace_listing_title"`
							CustomTitle                                string      `json:"custom_title"`
							CustomSubTitlesWithRenderingFlags          []struct {
								Subtitle string `json:"subtitle"`
							} `json:"custom_sub_titles_with_rendering_flags"`
							OriginGroup                             interface{}   `json:"origin_group"`
							PreRecordedVideos                       []interface{} `json:"pre_recorded_videos"`
							IsMarketplaceListingWithChildListings   string        `json:"__isMarketplaceListingWithChildListings"`
							ParentListing                           interface{}   `json:"parent_listing"`
							IsMarketplaceListingWithDeliveryOptions string        `json:"__isMarketplaceListingWithDeliveryOptions"`
							DeliveryTypes                           []string      `json:"delivery_types"`
							EstimatedDeliveryWindow                 interface{}   `json:"estimated_delivery_window"`
							HasFreeShipping                         bool          `json:"has_free_shipping"`
						} `json:"listing"`
						ID string `json:"id"`
					} `json:"node"`
					Cursor string `json:"cursor"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"end_cursor"`
					HasNextPage bool   `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"marketplace_feed_stories"`
			BuyLocation struct {
				BuyLocation struct {
					Location struct {
						ReverseGeocode struct {
							City string `json:"city"`
						} `json:"reverse_geocode"`
					} `json:"location"`
					ID string `json:"id"`
				} `json:"buy_location"`
			} `json:"buy_location"`
			MarketplaceSettings struct {
				CurrentMarketplace struct {
					Typename     string `json:"__typename"`
					IsMetricBase bool   `json:"is_metric_base"`
					ID           string `json:"id"`
				} `json:"current_marketplace"`
			} `json:"marketplace_settings"`
			MarketplaceSavedSearches struct {
				Edges []interface{} `json:"edges"`
			} `json:"marketplace_saved_searches"`
		} `json:"viewer"`
		MarketplaceSeoPage struct {
			Typename      string `json:"__typename"`
			SeoNavigation []struct {
				Typename              string `json:"__typename"`
				SeoURL                string `json:"seo_url"`
				SeoPageIsGeoAgnostic  bool   `json:"seo_page_is_geo_agnostic"`
				SeoLocalizedPageTitle string `json:"seo_localized_page_title"`
				VirtualCategory       struct {
					Name                        string `json:"name"`
					VirtualTaxonomyPublishState string `json:"virtual_taxonomy_publish_state"`
					ID                          string `json:"id"`
				} `json:"virtual_category"`
				ID string `json:"id"`
			} `json:"seo_navigation"`
			VirtualCategory struct {
				PivotsForLocation []struct {
					Name              string `json:"name"`
					IsGeoAgnosticOnly bool   `json:"is_geo_agnostic_only"`
					SeoInfo           struct {
						SeoURL string `json:"seo_url"`
						ID     string `json:"id"`
					} `json:"seo_info"`
					VirtualTaxonomyPublishState string `json:"virtual_taxonomy_publish_state"`
					ID                          string `json:"id"`
				} `json:"pivots_for_location"`
				ID string `json:"id"`
			} `json:"virtual_category"`
			ID string `json:"id"`
		} `json:"marketplace_seo_page"`
	} `json:"data"`
	Extensions struct {
		IsFinal bool `json:"is_final"`
	} `json:"extensions"`
}
