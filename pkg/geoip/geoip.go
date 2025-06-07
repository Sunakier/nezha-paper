package geoip

import (
	_ "embed"
	"net"
	"strings"
	"sync"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

//go:embed geoip.db
var db []byte

var (
	dbOnce = sync.OnceValues(func() (*maxminddb.Reader, error) {
		// Renamed 'db' to 'r' here to avoid conflict with package-level 'db' var
		r, err := maxminddb.FromBytes(db)
		if err != nil {
			return nil, err
		}
		return r, nil
	})
)

type IPInfo struct {
	Country       interface{} `maxminddb:"country"` // Changed to interface{}
	CountryName   string      `maxminddb:"country_name"`
	Continent     string      `maxminddb:"continent"`
	ContinentName string      `maxminddb:"continent_name"`
}

// IPInfoAlt is for cases where 'continent' is a map.
// Country is also changed to interface{} for consistency and to handle cases where both might be complex.
type IPInfoAlt struct {
	Country       interface{}            `maxminddb:"country"` // Changed to interface{}
	CountryName   string                 `maxminddb:"country_name"`
	Continent     map[string]interface{} `maxminddb:"continent"`
	ContinentName string                 `maxminddb:"continent_name"`
}

// Helper function to extract country code from interface{}
func getCountryCode(data interface{}) string {
	if data == nil {
		return ""
	}
	switch v := data.(type) {
	case string:
		return v
	case map[string]interface{}:
		// Common key for country code in MaxMind DB when country is a map
		if isoCode, ok := v["iso_code"].(string); ok {
			return isoCode
		}
		// Fallback: check if "code" key exists (less common but possible)
		if code, ok := v["code"].(string); ok {
			return code
		}
	}
	return ""
}

func Lookup(ip net.IP) (string, error) {
	dbReader, err := dbOnce() // Use dbReader (renamed from db to avoid confusion)
	if err != nil {
		return "", err
	}

	var record IPInfo
	err = dbReader.Lookup(ip, &record)

	if err != nil {
		// Check if the error is about 'continent' being a map (original fallback logic)
		if strings.Contains(err.Error(), "decoding value for continent") {
			var altRecord IPInfoAlt
			errAlt := dbReader.Lookup(ip, &altRecord)
			if errAlt != nil {
				return "", errAlt // Return error from this specific lookup
			}

			countryCode := getCountryCode(altRecord.Country) // altRecord.Country is now interface{}
			if countryCode != "" {
				return strings.ToLower(countryCode), nil
			}
			// If country code from altRecord is empty, original logic for IPInfoAlt was to return "", nil
			return "", nil
		}
		// If the error was not about 'continent' being a map, it might be a different unmarshalling issue
		// or a general lookup problem. The change of IPInfo.Country to interface{} should have prevented
		// the specific "cannot unmarshal map into type string" for the country field.
		return "", err // Return the original error if it's not the handled continent map case
	}

	// Successfully looked up with IPInfo (where Country is now interface{})
	countryCode := getCountryCode(record.Country)
	if countryCode != "" {
		return strings.ToLower(countryCode), nil
	}

	// Fallback to continent from standard record if country is empty
	// record.Continent is still string. If it can also be a map, IPInfo.Continent would need similar changes.
	if record.Continent != "" { // This is IPInfo.Continent, which is string
		return strings.ToLower(record.Continent), nil
	}

	return "", nil // Return empty string without error if no data found
}
