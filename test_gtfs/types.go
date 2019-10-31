package main

// GtfsRouteType -
type GtfsRouteType struct {
	Type    int
	Label   string
	Comment string
}

var gtfsRouteTypes = []GtfsRouteType{
	{0, "Tram, Streetcar, Light rail", "Any light rail or street level system within a metropolitan area."},
	{1, "Subway, Metro", "Any underground rail system within a metropolitan area."},
	{2, "Rail", "Used for intercity or long-distance travel."},
	{3, "Bus", "Used for short- and long-distance bus routes."},
	{4, "Ferry", "Used for short- and long-distance boat service."},
	{5, "Cable car", " Used for street-level cable cars where the cable runs beneath the car."},
	{6, "Gondola, Suspended cable car", "Typically used for aerial cable cars where the car is suspended from the cable."},
	{7, "Funicular", "Any rail system designed for steep inclines."},
	{100, "Railway Service", ""},
	{101, "High Speed Rail Service", "TGV (FR), ICE (DE), Eurostar (GB)"},
	{102, "Long Distance Trains", "InterCity/EuroCity"},
	{103, "Inter Regional Rail Service", "InterRegio (DE), Cross County Rail (GB)"},
	{104, "Car Transport Rail Service", ""},
	{105, "Sleeper Rail Service", "GNER Sleeper (GB)"},
	{106, "Regional Rail Service", "TER (FR), Regionalzug (DE)"},
	{107, "Tourist Railway Service", "Romney, Hythe & Dymchurch (GB)"},
	{108, "Rail Shuttle (Within Complex)", "Gatwick Shuttle (GB), Sky Line (DE)"},
	{109, "Suburban Railway", "S-Bahn (DE), RER (FR), S-tog (Kopenhagen)"},
	{110, "Replacement Rail Service", ""},
	{111, "Special Rail Service", ""},
	{112, "Lorry Transport Rail Service", ""},
	{113, "All Rail Services", ""},
	{114, "Cross-Country Rail Service", ""},
	{115, "Vehicle Transport Rail Service", ""},
	{116, "Rack and Pinion Railway", "Rochers de Naye (CH), Dolderbahn (CH)"},
	{117, "Additional Rail Service", ""},
	{200, "Coach Service", ""},
	{201, "International Coach Service", "EuroLine, Touring"},
	{202, "National Coach Service", "National Express (GB)"},
	{203, "Shuttle Coach Service", "Roissy Bus (FR), Reading-Heathrow (GB)"},
	{204, "Regional Coach Service", ""},
	{205, "Special Coach Service", ""},
	{206, "Sightseeing Coach Service", ""},
	{207, "Tourist Coach Service", ""},
	{208, "Commuter Coach Service", ""},
	{209, "All Coach Services", ""},
	{400, "Urban Railway Service", ""},
	{401, "Metro Service", "Métro de Paris"},
	{402, "Underground Service", "London Underground, U-Bahn"},
	{403, "Urban Railway Service", ""},
	{404, "All Urban Railway Services", ""},
	{405, "Monorail", ""},
	{700, "Bus Service", ""},
	{701, "Regional Bus Service", "Eastbourne-Maidstone (GB)"},
	{702, "Express Bus Service", "X19 Wokingham-Heathrow (GB)"},
	{703, "Stopping Bus Service", "38 London: Clapton Pond-Victoria (GB)"},
	{704, "Local Bus Service", ""},
	{705, "Night Bus Service", "N prefixed buses in London (GB)"},
	{706, "Post Bus Service", "Maidstone P4 (GB)"},
	{707, "Special Needs Bus", ""},
	{708, "Mobility Bus Service", ""},
	{709, "Mobility Bus for Registered Disabled", ""},
	{710, "Sightseeing Bus", ""},
	{711, "Shuttle Bus", "747 Heathrow-Gatwick Airport Service (GB)"},
	{712, "School Bus", ""},
	{713, "School and Public Service Bus", ""},
	{714, "Rail Replacement Bus Service", ""},
	{715, "Demand and Response Bus Service", ""},
	{716, "All Bus Services", ""},
	{717, "Share Taxi Service", "Marshrutka (RU), dolmuş (TR)"},
	{800, "Trolleybus Service", ""},
	{900, "Tram Service", ""},
	{901, "City Tram Service", ""},
	{902, "Local Tram Service", "Munich (DE), Brussels (BE), Croydon (GB)"},
	{903, "Regional Tram Service", ""},
	{904, "Sightseeing Tram Service", "Blackpool Seafront (GB)"},
	{905, "Shuttle Tram Service", ""},
	{906, "All Tram Services", ""},
	{907, "Cable Tram", "Cable car (San Francisco, US)"},
	{1000, "Water Transport Service", ""},
	{1100, "Air Service", ""},
	{1200, "Ferry Service", "Please use 1000"},
	{1300, "Aerial Lift Service", "Telefèric de Montjuïc (ES), Saleve (CH), Roosevelt Island Tramway (US)"},
	{1400, "Funicular Service", "Rigiblick (Zürich, CH)"},
	{1500, "Taxi Service", ""},
	{1501, "Communal Taxi Service", "Please use 717"},
	{1502, "Water Taxi Service", ""},
	{1503, "Rail Taxi Service", ""},
	{1504, "Bike Taxi Service", ""},
	{1505, "Licensed Taxi Service", ""},
	{1506, "Private Hire Service Vehicle", ""},
	{1507, "All Taxi Services", ""},
	{1700, "Miscellaneous Service", ""},
}
