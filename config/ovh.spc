connection "ovh" {
    plugin = "francois2metz/ovh"

    # Go to https://www.ovh.com/auth/api/createToken to create your application key,
    # secret and the consumer key
    # For the rights, GET with the path *
    # application_key = "CitIbyantOosuzFu"
    # application_secret = "phoagDakOywytMibfetJidloidvuenVo"
    # consumer_key = "einbycsAnmachCeOkvabicdifAdofdon"

    # OVH Endpoint
    # 'ovh-eu' for OVH Europe API
    # 'ovh-us' for OVH US API
    # 'ovh-ca' for OVH Canada API
    # 'soyoustart-eu' for So you Start Europe API
    # 'soyoustart-ca' for So you Start Canada API
    # 'kimsufi-eu' for Kimsufi Europe API
    # 'kimsufi-ca' for Kimsufi Canada API
    endpoint = "ovh-eu"

    # List of regions to query. Supports wildcards.
    # Defaults to all regions ("*") if not specified.
    # Examples:
    # regions = ["GRA"]                     # Single region
    # regions = ["GRA", "SBG", "BHS"]       # Multiple specific regions
    # regions = ["GRA*"]                    # All GRA regions (GRA, GRA9, etc.)
    # regions = ["GRA", "SBG", "BHS"]       # Multiple location wildcards
    # regions = ["*"]                       # All regions (default)
}
