//pagination
{
    "_links": {
        "self": {
            "href": "http://example.org/api/user?page=3"
        },
        "first": {
            "href": "http://example.org/api/user"
        },
        "prev": {
            "href": "http://example.org/api/user?page=2"
        },
        "next": {
            "href": "http://example.org/api/user?page=4"
        },
        "last": {
            "href": "http://example.org/api/user?page=133"
        }
    }
    "count": 3,
    "total": 498,
    "_embedded": {
        "users": [
            {
                "_links": {
                    "self": {
                        "href": "http://example.org/api/user/mwop"
                    }
                },
                "id": "mwop",
                "name": "Matthew Weier O'Phinney"
            },
            {
                "_links": {
                    "self": {
                        "href": "http://example.org/api/user/mac_nibblet"
                    }
                },
                "id": "windows_desktop",
                "name": "Sushritha Kolla"
            },
            {
                "_links": {
                    "self": {
                        "href": "http://example.org/api/user/spiffyjr"
                    }
                },
                "id": "spiffyjr",
                "name": "Kyle Spraggs"
            }
        ]
    }
}
