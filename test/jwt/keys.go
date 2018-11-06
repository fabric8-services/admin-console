package jwt

const (
	// PrivateKey1 the private key #1
	PrivateKey1 = `-----BEGIN RSA PRIVATE KEY-----
MIIEjgIBAAKB/g6X0vYyRg1xVKo8OUACpvDOXIqgQ3pUa5m8uVT2j4uz/Lb8kt2A
SsiwXwAFTpL6BcEgXM/eMrhodqvLT/CyDPEC24sXmgeTklAvs/Tqvm0YVJSRLWdX
6mh88BgtK7uq63fKxGUYo9YmlT7F9ejJtdkTaQApo8oa/zxHC3i0APSlaY0zU1F9
59QYStFVCK8Skn3cbODF/+SN1oxj3//t69sGnMWqfTllFcG7hD361wGp/PXI/n8U
1sVOe3rZEecHWwJF5kltLu6y7XyqHpPYrk8Qp+Fo1o3WfEcaG7P2JZVve35P+xYW
Bvm7ERsbqih5aSFJz8C3hKlwSQOEVtttAgMBAAECgf1c27yM4VriL0WP+6hQqI+h
wYEcnLDEumv22fB3tHe3gJeXzZq93p4AbEwW1a4nks8LG+N61W3qAtEgW5tTAalX
9dcOPiDkFSXzGZkD4Lnbefa7aRKBh+0S9fDT5ptikznGC32r0B65lMobp5IiuWds
6RY88rpKU3/PEETuzHtFO1kT/Dg/mny5FOTsmJjt+qjwROFCgiQBo/ZN33YTvF9v
aOcg5r2T3NmsDoXjorSHIFQ7+GC0Fj/rcUzLss9xgIsC93hJC1vt4+OlJSM4rIaW
GVI3fyVTZCjDqi2csXZCytTCSJRs1TggOo9jaine/xRi4r9jdXhmmyK9cXqhAn8+
AzcItixuziv7evlsHUsTTV/RR5TA9K+yXNm+9elX9NsaSVj0kRQeqTGWBua1bn0l
m+pvlsYimupSpJ3gmcrbcHUUdqvj655LIinWrLlN3GMLQyCvljk9BLWIrFQ6godR
rWEq2Kdni0vlm8f3h9hdR2XgS+dOvUKzE9R16rnJAn88PjbXRtTwf0hRXFYGigf5
DGYHt4WiIemduEC8Na+e2g6Sigdg7IqzQ04oSaWB+bM3tinyX+qIemaVL01VkJp4
7DMQrDx+KTsmmtYzBu7YLahZA0TqsnNRGeG+O3htDKsSKVUSsczoXu1Ar1+3nNim
RppCNMJGisNKFqRbjKaFAn8BEu8uEHGejaWHWm7dZ3h4YhuptTKnUNWGIkOHIh0j
b9MnlmObALQ3f7ijH4V5WOuD7jpWKmdODB7IxZ8SV7eCq2TrsM5zSQ5ZwMK2vBEN
fyab+FKll9Vv8BfwwQNIbCBJ0tXe9xeXHHt5A4SoDcs6elUSWF4uJ+ryzQId9K25
An8KMIe8H+nyh8TmphSS5JP2pwc29O6wfsXx/HFOpFIBL2bZmGkpFrlbGt5EaDiL
ZH3QxYoQyfJ0hSeGwkp1V5EZNPJqNofA2x57KCNk3B5YCFj6PVhRzj89D4CkWZDD
+SmSV9Vg5RwAjdXZZBBvkSL/9N8wpZXasqvXgz7nkUG1An8vOMeAmGunp9Be8r8Y
9QVTkvFtd44pLB/43sL6ZoXK0oPKgJUmCqPAYaJFRlfcdoHTQgymk0DuOKFhb4Ai
pfpBDkEIDzPInc3pWAtaqkZD7BdMnku3m78RNAOO97tMJPYQxxnNX16BqWm/SahW
chFXW3cvVxsaRJXr/a0gOYjI
-----END RSA PRIVATE KEY-----`

	// PublicKey1 the public key #1
	PublicKey1 = `-----BEGIN PUBLIC KEY-----
MIIBHjANBgkqhkiG9w0BAQEFAAOCAQsAMIIBBgKB/g6X0vYyRg1xVKo8OUACpvDO
XIqgQ3pUa5m8uVT2j4uz/Lb8kt2ASsiwXwAFTpL6BcEgXM/eMrhodqvLT/CyDPEC
24sXmgeTklAvs/Tqvm0YVJSRLWdX6mh88BgtK7uq63fKxGUYo9YmlT7F9ejJtdkT
aQApo8oa/zxHC3i0APSlaY0zU1F959QYStFVCK8Skn3cbODF/+SN1oxj3//t69sG
nMWqfTllFcG7hD361wGp/PXI/n8U1sVOe3rZEecHWwJF5kltLu6y7XyqHpPYrk8Q
p+Fo1o3WfEcaG7P2JZVve35P+xYWBvm7ERsbqih5aSFJz8C3hKlwSQOEVtttAgMB
AAE=
-----END PUBLIC KEY-----`

	// PrivateKey2 the private key #2
	PrivateKey2 = `-----BEGIN RSA PRIVATE KEY-----
MIIEjgIBAAKB/gpylWQfm3dy7OKY34RMCDNR89Jo6o1LVX0LNH6VK76WsTdkxass
nYY1m1hF3vWxzu1psJJthzIL5OnJI+Op3C7KFVYfPnLr8iLSMdb1bguPCuKRS6/l
stQpz8JhHY6RKfy5EG83OILQ+KsRYg+UAB+KkVqr8sqYEotWXC4hH57JIABfmIBM
4AxV5OUiMllOJWcYHVakOfS1TE1W0vtBReKFB23wxw/C57WkwXLquvVcrEZqxWxV
Y+wNTnYh8XlkY1jX7OW9QGSHjZK18PjxurxwvOJ5hF3xR2X1q7ffn3FPCeUt5zDS
E/V9q0kVK17u3QfMxhDqLYR5b/0uZJx3AgMBAAECgf4E/T1MaC+1NmPbvoeXRTvx
TiSzCblhKmWz5mL2RER0qsAMpQokuZSsX/NEj3FvQa+A/yT8eGPEyZtS7eQ+t4JX
sdfInfkTpoumh1yXu/MGgBQBqMNNR9NDsIfv2rLjv30enD417lgFWMg34YBD0jjQ
1zqc41p512+brO0udlEEML0g/uvxVwHH+vkPW/UhQcQaAFSkKnjV23ZZOKwexqvZ
r4aOM9rNC9QvaEVEPkQS0jbnyn8QEfyfEW/e09u4wbtnLJanERG/q0hQwnf6jV03
uZCWjbrqO2lxYPX9Eh7z8YAHuV1DqIBEUxo/eW++XBFqi6NDfaIg2xfq2GFFaQJ/
NQYOcbvin8rdlB7OYHLLKvQOBLRQ2yZpSbzxsNyhysA/GsIK16o8LXE4Yx5ajRt0
DD7ucjvV5kyYSrY/lxp3wB9Sxe1GH5bGKac2233CghhN+aE4RXCX7Ys2vb9r054M
OVykhd7L7OQlLCOiz1XEWEjTSOsQnDnnfkLVcbdXKwJ/MnD6rD6n2AiJdJb+42ly
H3OjXDytKFIbwbT2cNL+PqQx3IoZShFswG/k/ur/qrbRChFOPe1Ra+MfuG1/G04h
aGv3SA8/4gnF6lZSeRTA2aD/1BWJjCcDMBikKyCB/9FIBtSxd7Dc7RKWz+r6ofG2
AhH2hUvyBgnwlYNFhmlp5QJ/NKgGEiaPATcaYu1Q7/EwGEDz5vIW1fvIVaUgtA33
Un0mbfFDiTcSueIVKEHNlRItZbfdXm5Tlnh8SL3CWtG7GH1C2zIuEFLQCi93x/OV
BVMTpJLZagRNrGYy/66oaygqZZC+Bf/iridHTxU4qwQ2j6IKyQA/p5XNcdL3Ww3C
rwJ+YukopJo4h2g2Imn1Z/tdzk70B+rzoh1gUHiKyLL12+Ad5ljnPlbs6m6AnWAz
+I6FWziUNNsozmaRzRTqKqEK4bjVLni4zIZdkyeykbwgbqFHCJszHaFinu2y/t7A
DORWWQ668fnfPwM0uIIO94EDZwvSXZCPB0h2fLjtnKUZAn8X726zlS98DDO7RGqp
/Jw3/QlyCUermTnKLeHwgtuAKeAjJOqF/vFnf3mweJ5KHAHheLfOnOIl+58UVjj8
JCz5cPdVRH8DBD60gm3pymHTxMnq6EjJ99ZWmY2nHiNRN/rrxTPQKhQMHCippPjM
M4w5F9oxOY13Gw0MWXn9mrY9
-----END RSA PRIVATE KEY-----`

	// PublicKey2 the public key #2
	PublicKey2 = `-----BEGIN PUBLIC KEY-----
MIIBHjANBgkqhkiG9w0BAQEFAAOCAQsAMIIBBgKB/gpylWQfm3dy7OKY34RMCDNR
89Jo6o1LVX0LNH6VK76WsTdkxassnYY1m1hF3vWxzu1psJJthzIL5OnJI+Op3C7K
FVYfPnLr8iLSMdb1bguPCuKRS6/lstQpz8JhHY6RKfy5EG83OILQ+KsRYg+UAB+K
kVqr8sqYEotWXC4hH57JIABfmIBM4AxV5OUiMllOJWcYHVakOfS1TE1W0vtBReKF
B23wxw/C57WkwXLquvVcrEZqxWxVY+wNTnYh8XlkY1jX7OW9QGSHjZK18Pjxurxw
vOJ5hF3xR2X1q7ffn3FPCeUt5zDSE/V9q0kVK17u3QfMxhDqLYR5b/0uZJx3AgMB
AAE=
-----END PUBLIC KEY-----`
)