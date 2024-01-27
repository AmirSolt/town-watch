package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
)

type ArcgisResponse struct {
	Features []ArcgisReport `json:"features"`
}

type ArcgisReport struct {
	Attributes ArcgisAttributes `json:"attributes"`
	Geometry   ArcgisGeometry   `json:"geometry"`
}

type ArcgisGeometry struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type ArcgisAttributes struct {
	ObjectId         string    `json:"OBJECTID"`
	OccDateEst       time.Time `json:"OCC_DATE_EST"`
	OccDateAgol      time.Time `json:"OCC_DATE_AGOL"`
	ReportDateEst    time.Time `json:"REPORT_DATE_EST"`
	ReportDateAgol   time.Time `json:"REPORT_DATE_AGOL"`
	EventUniqueId    string    `json:"EVENT_UNIQUE_ID"`
	Division         string    `json:"DIVISION"`
	PremisesType     string    `json:"PREMISES_TYPE"`
	Hour             int16     `json:"HOUR"`
	CrimeType        string    `json:"CRIME_TYPE"`
	Hood158          string    `json:"HOOD_158"`
	Neighbourhood158 string    `json:"NEIGHBOURHOOD_158"`
	Hood140          string    `json:"HOOD_140"`
	Neighbourhood140 string    `json:"NEIGHBOURHOOD_140"`
	Count            int16     `json:"COUNT_"`
	LongWgs84        float32   `json:"LONG_WGS84"`
	LatWgs84         float32   `json:"LAT_WGS84"`
	LocationCategory string    `json:"LOCATION_CATEGORY"`
}

func (server *Server) FetchReports(fromDate time.Time, toDate time.Time) (*ArcgisResponse, error) {
	toDateStr := fmt.Sprintf("AND OccDateAgol <= date'${toDate} 00:00:00'")
	where := fmt.Sprintf("OccDateAgol >= date '${fromDateStr} 00:00:00' %s", toDateStr)
	endpoint := fmt.Sprintf("https://services.arcgis.com/S9th0jAJ7bqgIRjw/ArcGIS/rest/services/YTD_CRIME_WM/FeatureServer/0/query?where=%s&objectIds=&time=&geometry=&geometryType=esriGeometryEnvelope&inSR=&spatialRel=esriSpatialRelIntersects&resultType=none&distance=0.0&units=esriSRUnit_Meter&relationParam=&returnGeodetic=false&outFields=*&returnGeometry=true&featureEncoding=esriDefault&multipatchOption=xyFootprint&maxAllowableOffset=&geometryPrecision=&outSR=&defaultSR=&datumTransformation=&applyVCSProjection=false&returnIdsOnly=false&returnUniqueIdsOnly=false&returnCountOnly=false&returnExtentOnly=false&returnQueryGeometry=false&returnDistinctValues=false&cacheHint=false&orderByFields=&groupByFieldsForStatistics=&outStatistics=&having=&resultOffset=&resultRecordCount=&returnZ=false&returnM=false&returnExceededLimitFeatures=true&quantizationParameters=&sqlFormat=none&f=pjson&token=", url.QueryEscape(where))
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making request to Arcgis API: %w", err)
	}
	defer resp.Body.Close()

	var response ArcgisResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding ArcgisResponse: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	vErr := validate.Struct(response)
	if vErr != nil {
		log.Fatal("Error a variable is missing from .env")
	}
	return &response, nil
}

// let rawReports: any[] = []

// const DATE_TYPE = "REPORT_DATE_AGOL"

// const fromDateStr = UTCToStr(fromDate)
// const query_date_str = `${DATE_TYPE} >= date '${fromDateStr} 00:00:00'`
// const whereStatementUrlified = encodeURIComponent(query_date_str)
// const rawReportsUrl = `https://services.arcgis.com/S9th0jAJ7bqgIRjw/ArcGIS/rest/services/YTD_CRIME_WM/FeatureServer/0/query?where=${whereStatementUrlified}&objectIds=&time=&geometry=&geometryType=esriGeometryEnvelope&inSR=&spatialRel=esriSpatialRelIntersects&resultType=none&distance=0.0&units=esriSRUnit_Meter&relationParam=&returnGeodetic=false&outFields=*&returnGeometry=true&featureEncoding=esriDefault&multipatchOption=xyFootprint&maxAllowableOffset=&geometryPrecision=&outSR=&defaultSR=&datumTransformation=&applyVCSProjection=false&returnIdsOnly=false&returnUniqueIdsOnly=false&returnCountOnly=false&returnExtentOnly=false&returnQueryGeometry=false&returnDistinctValues=false&cacheHint=false&orderByFields=&groupByFieldsForStatistics=&outStatistics=&having=&resultOffset=&resultRecordCount=&returnZ=false&returnM=false&returnExceededLimitFeatures=true&quantizationParameters=&sqlFormat=none&f=pjson&token=`
// const response = await fetch(rawReportsUrl)
// let data = await response.json() as any
// if (response.ok) {
// 	rawReports = data["features"] as any[]
// } else {
// 	throw fastify.httpErrors.conflict(JSON.stringify(data))
// }

// // filter null fields
// rawReports = rawReports.filter(raw => {
// 	return (
// 		raw["attributes"] != null ||
// 		raw["attributes"]["CRIME_TYPE"] != null ||
// 		raw["attributes"][DATE_TYPE] != null ||
// 		raw["attributes"]["HOUR"] != null ||
// 		raw["geometry"] != null
// 	)
// })

// function getInsertValuesStr(rawReports: any[]) {
// 	return rawReports.map(raw => {
// 		const date_at = setDateToSpecificHour(raw["attributes"][DATE_TYPE], raw["attributes"]["HOUR"])

// 		return `
// 			(
// 				${date_at},
// 				${region}_${raw["attributes"]["EVENT_UNIQUE_ID"]},
// 				${removeNeighExtraChars(raw["attributes"]["NEIGHBOURHOOD_158"])},
// 				${raw["attributes"]["LOCATION_CATEGORY"]},
// 				${crimeTypeCleaning(raw["attributes"]["CRIME_TYPE"]) as CrimeType},
// 				${region},
// 				ST_Point(${raw["geometry"]["x"]}, ${raw["geometry"]["y"]}, 3857),
// 				${raw["geometry"]["x"]},
// 				${raw["geometry"]["y"]}
// 			)
// 	`}).join(",")
// }

// return await fastify.prisma.$executeRaw`
// 	INSERT INTO "Report" (
// 			reported_at,
// 			external_src_id,
// 			neighborhood,
// 			location_type,
// 			crime_type,
// 			region,
// 			point,
// 			lat,
// 			long
// 		)

// 		VALUES
// 		${getInsertValuesStr(rawReports)}
// 		ON CONFLICT DO NOTHING;
// `

// }

// function UTCToStr(date: Date): string {
// // Format the date and time
// let year = date.getUTCFullYear();
// let month = ("0" + (date.getUTCMonth() + 1)).slice(-2);
// let day = ("0" + date.getUTCDate()).slice(-2);
// let hours = ("0" + date.getUTCHours()).slice(-2);
// let minutes = ("0" + date.getUTCMinutes()).slice(-2);
// let seconds = ("0" + date.getUTCSeconds()).slice(-2);

// return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
// }

// function removeNeighExtraChars(inputString: string): string {
// let result: string = "";
// let flag: boolean = false;   // True when between "(" and ")"

// for (let char of inputString) {
// 	if (char === "(") {
// 		flag = true;
// 	} else if (char === ")") {
// 		flag = false;
// 	} else if (!flag) {
// 		result += char;
// 	}
// }

// return result;
// }

// function setDateToSpecificHour(epochTime: number, hour: number): Date {
// const date = new Date(epochTime);
// date.setUTCHours(hour, 0, 0, 0);
// return date;
// }

// function crimeTypeCleaning(rawCrimeType: string) {
// return rawCrimeType.toUpperCase().replace(/\s+/g, '_');
// }
