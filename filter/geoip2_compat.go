package filter

import (
	"fmt"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"strconv"
)

type GeoIP2LegacyCompatDecoder struct {
	GeoIPDecoder
}

func (g *GeoIP2LegacyCompatDecoder) decode(ip string, field string, doc map[string]interface{}) {
	if v, ok := doc[fmt.Sprintf("%s.geoip2.country.iso_code", field)]; ok {
		doc[fmt.Sprintf("%s.country_code", field)] = v
	}
	// According to https://dev.maxmind.com/geoip/geoip2/whats-new-in-geoip2/#Custom_Country_Codes,
	// "For IP addresses where we used A1 or A2, we now set the appropriate key in the traits object.
	// The A1 code corresponds to the is_anonymous_proxy key, while A2 corresponds to is_satellite_provider."
	//
	// The country code field in geoip legacy was overridden in cases to denote whether an IP is an
	// anon proxy, or if it is a satellite provider. In geoip2, maxmind separated the proxy/satellite
	// demarcation into a separate 'traits' field.
	//
	// Thus, for compatibility reasons, we will set country_code appropriately if either is_anonymous_proxy
	// or is_satellite_provider are set.
	if is_anon_proxy_iface, ok := doc[fmt.Sprintf("%s.geoip2.traits.is_anon_proxy", field)]; ok {
		if is_anon_proxy_str, ok := is_anon_proxy_iface.(string); ok {
			if val, err := strconv.ParseBool(is_anon_proxy_str); val && err == nil {
				doc[fmt.Sprintf("%s.country_code", field)] = "A1"
			}
		}
	} else if is_satellite_iface, ok := doc[fmt.Sprintf("%s.geoip2.traits.is_satellite", field)]; ok {
		if is_satellite_str, ok := is_satellite_iface.(string); ok {
			if val, err := strconv.ParseBool(is_satellite_str); val && err == nil {
				doc[fmt.Sprintf("%s.country_code", field)] = "A2"
			}
		}
	}

	// No mapping for country_code3 from geoip2
	// doc[fmt.Sprintf("%s.country_code3", field)] = record.CountryCode3

	if v, ok := doc[fmt.Sprintf("%s.geoip2.country.name", field)]; ok {
		doc[fmt.Sprintf("%s.country_name", field)] = v
	}

	if v, ok := doc[fmt.Sprintf("%s.geoip2.location.metro_code", field)]; ok && v.(uint) > 0 {
		doc[fmt.Sprintf("%s.dma_code", field)] = v
	}

	// No mapping for area code from geoip2
	// https://dev.maxmind.com/geoip/geoip2/whats-new-in-geoip2/#Area_Code
	// doc[fmt.Sprintf("%s.area_code", field)] = record.AreaCode

	if v, ok := doc[fmt.Sprintf("%s.geoip2.postal.code", field)]; ok && v != "" {
		doc[fmt.Sprintf("%s.postal_code", field)] = v
	}

	if v, ok := doc[fmt.Sprintf("%s.geoip2.location.latitude", field)]; ok {
		doc[fmt.Sprintf("%s.latitude", field)] = v
	}

	if v, ok := doc[fmt.Sprintf("%s.geoip2.location.longitude", field)]; ok {
		doc[fmt.Sprintf("%s.longitude", field)] = v
	}

	if v, ok := doc[fmt.Sprintf("%s.geoip2.city.name", field)]; ok && v != "" {
		doc[fmt.Sprintf("%s.city", field)] = v
	}

	if num_subdivisions_iface, ok := doc[fmt.Sprintf("%s.geoip2.subdivisions.length", field)]; ok {
		if num_subdivisions, ok := num_subdivisions_iface.(int); ok {
			// we don't really care to loop over the number of subdivisions.
			// we only care that there's at least one, and if there is, we'll
			// take the first one.
			if num_subdivisions > 0 {
				if v, ok := doc[fmt.Sprintf("%s.geoip2.subdivisions.0.iso_code", field)]; ok {
					doc[fmt.Sprintf("%s.region", field)] = v
				}

				if v, ok := doc[fmt.Sprintf("%s.geoip2.subdivisions.0.name", field)]; ok {
					doc[fmt.Sprintf("%s.region_name", field)] = v
				}
			}
		}
	}
}

func init() {
	factory.RegisterFilter("geo_ip2_legacy_compat", func(args []string) (filterGeoIPCity api.Filter, err error) {
		decoder := new(GeoIP2LegacyCompatDecoder)
		return NewFilterGeoIP(args, decoder)
	})
}
