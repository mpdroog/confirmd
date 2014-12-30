confirmd
========

Proxy public webservices through one simple JSON-URL.

Supported lookups:
* EU VIES (European VAT-number lookup)
* OpenKvK (Dutch Chamber of Commerce lookup)
* Postcode (Dutch postal code lookup)

Why this abstraction?
* Possible to cache requests
* Indirection, making replacement a breeze

TODO:
* Cache valid responses
* Postal check for BE, FR, UK?