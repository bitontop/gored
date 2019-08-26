package conf

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"github.com/bitontop/gored/exchange"
)

func Exchange(name exchange.ExchangeName, config *exchange.Config) {
	config.ExName = name
	switch name {
	case exchange.BINANCE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITTREX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.STEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITMEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.KUCOIN:
		config.API_KEY = ""
		config.API_SECRET = ""
		config.Passphrase = ""
		break

	case exchange.BITMAX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITSTAMP:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.OTCBTC:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.HUOBI:
		config.API_KEY = "955276ce-3738027c-2669c1ba-520b6"
		config.API_SECRET = "359c4c74-6dc2dfe6-628ad35d-1032a"
		break

	case exchange.BIBOX:
		config.API_KEY = "b5b45115948fff3db063e84924537cbe13cb9794"
		config.API_SECRET = "536be20857fc6295aa56a26607fa852087378af2"
		break

	case exchange.OKEX:
		config.API_KEY = "bf1a3702-50d8-4efc-b631-739b84966b5a"
		config.API_SECRET = "6630E22F16A786E3F5525D54E1A4210A"
		config.Passphrase = "tanboo"
		config.TradePassword = ""
		break

	case exchange.BITZ:
		config.API_KEY = ""
		config.API_SECRET = ""
		config.TradePassword = ""
		break

	case exchange.HITBTC:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.DRAGONEX:
		config.API_KEY = "824ec91f624c5f0287303116f3f5ba33"
		config.API_SECRET = "54b4a115637b5727af9d4e7e7a35ef97"
		break

	case exchange.BIGONE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITFINEX:
		config.API_KEY = "6nSXMTuh2XvQS7HCCFSS9koKNR7wHW9462msqwxwSWJ"
		config.API_SECRET = "lnHfFCAppXSDYg8P7TuZBVB6oyrdRFMc0wCwxE4o6hm"
		break

	case exchange.GATEIO:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.IDEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.LIQUID:
		config.API_KEY = "918463"
		config.API_SECRET = "RsUNgGO6EtV4Fo4ATLwl8m3eDpRe8piMnKBQl3Y8hCR67PQu8si+5J/6rlKqV1/fUHYfAk/JuLO9ez7xJmKsjg=="
		break

	case exchange.BITFOREX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TOKOK:
		config.API_KEY = "5ecac940-b495-48c6-9bd1-dd0e1c6b6b95"
		config.API_SECRET = "14a9a62e-fcaf-49a8-a4a9-dbaf5921bbff"
		break

	case exchange.MXC:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITRUE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TRADESATOSHI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.KRAKEN:
		config.API_KEY = ""
		config.API_SECRET = ""
		config.Two_Factor = ""
		break

	case exchange.POLONIEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINEAL:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TRADEOGRE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINBENE:
		config.API_KEY = "1"
		config.API_SECRET = "2"
		break

	case exchange.IBANKDIGITAL:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.LBANK:
		config.API_KEY = "c916eb09-f6a3-47e5-ad24-b995bc3320f3"
		config.API_SECRET = "MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBANofzf//Pl3/etW/Z4L5WL+Nbxl+Ac2KAwWI8qb6ciYygWl+idlYIQn/1yUvwAwo5SyYr0sqLQaMUL4C+nCGYjwe2FUIRTViKNycoZM94FxkuvWsVlpToDmzMOhXKvKU7+Yj/quYfmjqkMR7vG0VuPP1CYyrR2epC709tvPzuUVzAgMBAAECgYEAsH6NrDe3Gk4P8Ya31iW2pwBlRkZMZSjoOwFN/siltrylNFxcZE5IJZQrXP6fMfehQI2nQXW2CxdcefNk+8nxD2AG8FxuLSnYHdksKTTfDCjG7ynYxX1TdDrWS7k7UKITRxF1FgQgXN2Taj8YDVs8eJlzG1H6hY/4WMXmhqtu2dECQQDyfhLy80DOX1YWMHWnWD+aqs8lk7OZjugrnp66mYGiMBWGd7vm3ztlxjRrIBTMRvYHSdhvj+FGFJIRr7+MTCRpAkEA5kY9wsaxN+rqDuH90B0KycmXYsftGPaIVdcgK1ypi2WRa1Cl87f/C1G0ukc7qHxszuW8j3TPemTe/foKxECvewJBALjf77h0RqtQOgTOy1RbTpqvsSBX/GyNbGqdEyz2jcPGXxLWxFYfSVytgRdPLSwUycwCu9VKX5ibZEXBtQrUnkkCQAaA91+X8wtsRB4VffCx0UsvuWnd0bGBzQn3oH35CQTZ4oiQZ4+Bo99y+FLGjkXM9dnGHVRD7VQ8oxuzQziVxx8CQQCzcU4EKeC4YRDxc3AY9CtW9Yfz4Pzrd7WeK/twUJu4sIeYMPb2dmb9wZuetNjESirXujoNcH7Brb4bhl8ImlHw"
		break

	case exchange.BITMART:
		config.API_KEY = "3e43d35bbeefeb881b4de213dc01042f"
		config.API_SECRET = "fa65bc32e69c1b7f2f8965ca7cf5d4fa"
		config.Passphrase = "key3" // key name
		break

	case exchange.BIKI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITATM:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.DCOIN:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.GEMINI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINTIGER:
		config.API_KEY = "fe6d48e4-cf08-4566-b217-3ca8e8c4014c"
		config.API_SECRET = "NDRhYWViNmE3YzY2OGMzM2FhMjUyZjUxODdjNDBjNTE4NTczOTEyMzcxYmQ1MDg2OGUyZDI4ZGU2MGZhZWMzOA=="
		break

	case exchange.HUOBIDM:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BW:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITBAY:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.DERIBIT:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.OKEXDM:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.GOKO:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BCEX:
		config.API_KEY = "baf857298a8b06c66dea5213de37c076"
		config.API_SECRET = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDO+V75tu9aCotu
Ri53yQQLb9J3sUi3pVW0VKyVyAJlQ1g6BX7xLUe5/puN9Ee86CwTTKmLSTz7JhrB
R/ixUeIIrCa78bB6e8E4NL8qNPfcmRbDcxWQ1/4kA+y/MZJsmeNt2J/gmK8g4sL+
uZ3eR11knNmPfpNHU2YtuNExBZTAwriPy5d0l3+yWKuVG4EWzHGKD2YccvSZVIFO
wjr/LYWXOiYBAum5fzUX97gGVcycLoM7kY+ih2keh6upqt239YfeyMqZyY7+Rzhf
qY006bUinQS1ZECmjnHUmKwejZAPyrQ/31uJgkLukRXRHax2N/StivQDK8zxfa9w
Yg3qTuF1AgMBAAECggEBAKEwBi5VsJwgACx2TbQOAa9ie7epgqc7BL72/p17cZop
U2mEZDMxf12zkrN/3eqZqdGx74xBElPJfiauhVibG3yDjXrrI/SPso+yTHpzW+EZ
/GdklhQCkrK7t8HCunUHd95RSVmhryneT9wO9Ipqa6pymOCuw8ZVhgrvl0MlRI7E
71rLzkHs3UeFTQOjv/5vJF97JU/kzjitV0bNISHFgGRFDfH+hjUTYjAXnbtUhzWv
gf9nhonSBCXHMjwKN7MIM0bZOr/CtbKESzsRLGNjWf/EutURhLliaaYS3Zf+CdP6
ZibEOkAcfPpO0I70x5+Ho4wG79FohTLoOjm0uZJFzbUCgYEA+KrTgHzw/c1T7nG1
wY9qhH15y/tFtDtaXjxej28HNE8EHWTmwyUIp5hg9U8b42tkNNz/vQ5Q8vHjkBhF
1nl5GyVl4hosEhES6lavJ2P/Jt2F9SoYOgvxHqRohYbH4y02BWYKs6WqHJaUFTsK
2B0LlH9tlguDj9qjZ9WZ7FJt/0cCgYEA1RPOOlHkwdLeWG1H7Rx7DlNf2BRGjjB8
eiKId1lxyfzfNuzZ/JRx8SEaT2X8BMGWROKRK7pBLOuI21/6G79Pxxuu1DWyCQqB
qhKbOy8orJRgMVOMx8wE66kGhsqDeRonh3TJOtST3+iCPzglgVxfQVwhyrAD5NCG
PWH/mrGwD2MCgYEAzW9bongWJKgAYiqhDRMt3d1HxUSG1pp+UwIu4PLKEeYBsUMN
/kRXPRZ/a8p6cMzlEWNPCGKOb9d0uDPFZqYeblXcMQqMRDTE2sLYm4NaZUJ4DA5F
y5bYEgejrkSmWMGeMqGVz8ramhmwp0WK9PYx/fG0mFRU0YDApOTr8Dg9VbsCgYBR
+vfKqn2IMViIzyrwSJfz8BIdMdfflzodR7IXsVs9asR6/m/0ZSzdqG3WBJgNQGpP
gJh4KYYwAUM7nFa/XEEWi0kdrrccEWXICLae88sDc2b7M3kj2hQ+k17Gd55T9sMk
s8NElkt6x5ttNW4AsoiXvhnmQQiOfchYT58nZpwlnwKBgBnrISoODvISEz4lkNBk
euHoGaR0Sa5aWJHv+xOhhbksAuIXHDibfT4Pw2v6AxryWT1Kt2GLAJrJYFmEEe5R
M0xyX9XiuAcJAmkK6zb9YFpKrp+ASfI+ufjDe4KCSmrPd728mv/FgRJ8mMfmod5t
0QB6J7s9DGKyKNKVTfNHRBvV
-----END PRIVATE KEY-----
`
		break

	case exchange.DIGIFINEX:
		config.API_KEY = "15d276b8f96138"
		config.API_SECRET = "d3dce802d35e921b97232915f49f161905d276b8f"
		break

	case exchange.LATOKEN:
		config.API_KEY = "api-v1-a20f884758774b22341e4ba4b5a96498"
		config.API_SECRET = "api-v1-secret-a6acbdee371f31d51d0472a591a27cf5"
		break

	case exchange.VIRGOCX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.ABCC:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BYBIT:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.ZEBITEX:
		config.API_KEY = "PpouZdAESSwxUfjDZOk3jLsf8hXk2KW7BtPYdl8I"
		config.API_SECRET = "XxFWbUAZHHVROIjXMxxPY9xOIp0NGCmJLUJL2BIG"
	}
}
