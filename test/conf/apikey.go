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
		config.API_KEY = "CIzP4OFNi15loANfLehB1nIN4rKGZhKeirXXW9UcExMAvGDj7UdAqbHURK6Z5v0z"
		config.API_SECRET = "84TJ0diip2UfGDAQ6KzUCHP8psIS72ewE20MBrlXsbCmEJlsWMdEARkFlSuMDCoF"
		break

	case exchange.BITTREX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINEX:
		config.API_KEY = "A017C186951B452788029681854B651E"
		config.API_SECRET = "7EAF00A996AF4340B67B156F76DB56645E74935B4650B381"
		break

	case exchange.STEX:
		config.API_KEY = " "
		config.API_SECRET = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxIiwianRpIjoiN2U4ZGNjYTE0OGQ0MDMxMTEzZjFmMzFkNGJmYjA4OWJlYTk2YWYxZmFiN2VmMzgyZTI1ZTcwNDdkOWZkNThjODZjMGM3MTUxNmYwY2ZmNzgiLCJpYXQiOjE1NzMxNzU5NDksIm5iZiI6MTU3MzE3NTk0OSwiZXhwIjoxNjA0Nzk4MzQ5LCJzdWIiOiIzMjc5NjciLCJzY29wZXMiOlsicHJvZmlsZSIsInRyYWRlIiwid2l0aGRyYXdhbCIsInJlcG9ydHMiLCJwdXNoIiwic2V0dGluZ3MiXX0.eaVyJ18T0UnZkinI6LhJNkimk-Do-h1MRuZd6w8oXXYjZOpbuJ9aMZ7HFredhnjQ6RUhkJrG29T26jUaRDwrMIaUYnOa7_rFdMYXOuaqugXc7RWypK74I3rWPJQ4Ber5N4Ky26tQRyvoCnu-S1wpZZGmfRiPCL7cgQHiQOqR5p8E1hBfEFRAJ-G2oDdWnRwPZ7oGjojnqCM0c2anBBIAeSPcu_8vhGmG-Of85bulh-MX-1-LS7rTuh1Xr_MWd30XH-niCPHP0lxcT3ONLD6Qaadd-rmtqthGH6W5_Urw6l8RjwnSC5It4EnnBWzCQQfTE4lvT0TEJbooFNQ4qHFZ1XudTqKp0z1fiZyDLpn4PIV1spUJGTgDkshx1y1Ch-ppVGq0MEt5FlZRKaKHSvgIz0U9PTyeBne3pwBA-YBd85oSlNRvpPAdhLuN6Hjzy0NLNs_5zzpjUFz-oUsYSej97SUYQe8ZWnB6DNF4X715oR6bvzMyoVV7Z4qTyy0CHku9-I5HaUdqnP4CakF-t6AsASz4lJofz-nRVMAsuaRJzUQs8L8ouo5PGm491f27nEfYg7EJAIEQkwCIZFxZGENQ524hPPD5tsXv5mD7KEcqW3AZhRxDw_hghfO8w3WIHJklnIBySdXslo25pdScyPl3NRSrFJtBO9jL90n_sXZBQ4o" //"S2S Token"
		config.ExpireTS = 1573081822
		break

	case exchange.BITMEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.KUCOIN:
		config.API_KEY = "5dbc72e04c06872bb6883eba"                //"5cf9a1e9134ab70a891ee2f3"
		config.API_SECRET = "03202cda-8725-458e-a29a-f94bf5fa11c0" //"ed5f967f-de14-48ee-a9ac-2617da658868"
		config.Passphrase = "11223344"
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
		config.API_KEY = "78d5a4a3-4546d9dc-f391ad89-h6n2d4f5gh"
		config.API_SECRET = "73318fc9-b8323aa4-8babaae1-f1c2f"
		break

	case exchange.BIBOX:
		config.API_KEY = "b5b45115948fff3db063e84924537cbe13cb9794"
		config.API_SECRET = "536be20857fc6295aa56a26607fa852087378af2"
		break

	case exchange.OKEX:
		config.API_KEY = "233110f2-9a4c-4643-8c1b-29de3c28c79f" //"e2ebe992-53de-4d78-8e0e-26d7d55fae28"
		config.API_SECRET = "29D9A710B7E9539CB38AF344252055DF"  //"15B55E55207920DC4F777845D374C468"
		config.Passphrase = "Bitontop1"                         //"bitontop"
		config.TradePassword = "Bitontop1"
		break

	case exchange.BITZ:
		config.API_KEY = "98a8d929ecb7ff1ae6d5a8ea748ccfa3"
		config.API_SECRET = "IvMsC8hzckM465a7yucyDFtOYPTQK3qYlegWBJPJLGHwmAfRfSxGwz7VKj6F6R4Q"
		config.TradePassword = "Bitontop2"
		break

	case exchange.HITBTC:
		config.API_KEY = "EWOYXwQW_x6WCdVJmnBSvBmmjglDJ9g9"
		config.API_SECRET = "W3ZqUE38ZTNYx7DHeqe78UwP6u5lxtO1"
		break

	case exchange.DRAGONEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BIGONE:
		config.API_KEY = "9a86c254-7e7e-4aaa-97b6-1dc8d0cc46c2"
		config.API_SECRET = "F9A1CD65EC0C8589ECC0D77B3715D24D0DA4CAFC3C0AA539CA63FBDFB3F6E702"
		break

	case exchange.BITFINEX:
		config.API_KEY = "6nSXMTuh2XvQS7HCCFSS9koKNR7wHW9462msqwxwSWJ"
		config.API_SECRET = "lnHfFCAppXSDYg8P7TuZBVB6oyrdRFMc0wCwxE4o6hm"
		break

	case exchange.GATEIO:
		config.API_KEY = "3DF053C7-63F2-44BC-80ED-4F0C53362F95"
		config.API_SECRET = "d4dd52a544cc8543c74d99484f0c447bf7ade0753a2200e165c8925a230654b7"
		break

	case exchange.IDEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.LIQUID:
		config.API_KEY = "1093696"
		config.API_SECRET = "ujoyVTbYYrr5IHyWJMY3BdfgdCAL+oKUXaFpuWgkQoRSjh8vzPQ3z9LKfN1QHWBgiuhhaRkQtfTK7jgUSm5e0g=="
		break

	case exchange.BITFOREX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TOKOK:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.MXC:
		config.API_KEY = "mx0NLMLlQDrIdGUu39"
		config.API_SECRET = "fa4b0eaab51f4db1a778ab04afeb8f3a"
		break

	case exchange.BITRUE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TRADESATOSHI:
		config.API_KEY = "70eeba574f2341a8bb790377dab6458e"
		config.API_SECRET = "hVwNuRsxRiB0JNRCu1+zyi7xYacV/coU9sud9ZkTl4A="
		break

	case exchange.KRAKEN:
		config.API_KEY = "NFGFpUMAvncdDA/fRALhW3EwTgtSluVHdbImx6Oulvz4u/SFoGiCHaNa"
		config.API_SECRET = "TU4i7Xb5Av65NjgeVd1gR19lNlWnwBWnMywCixGwyLiKaTbYP+5WmU+zfExYsR42PqfVzK+HNSHnefq9we0fiQ=="
		config.Two_Factor = ""
		break

	case exchange.POLONIEX:
		config.API_KEY = "YPLOL03K-MAP6N9IG-ODFMW7IA-5YDDS2OG"
		config.API_SECRET = "a618af875106b4ab0c1c2c3d1b880c0e987011927bafa3aaeb9d3ef2d7fdacd07c146db8822639a9c2400621baffddc7b78358854caecdf462d08a98e621f34d"
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
		config.API_KEY = "891461707962c31758200cee1938a5b0"    //"2aa98cdba8f2f104973443ff27375f80"
		config.API_SECRET = "b079dd6190784b618bfd5e308ca0debb" //"6e8154b24d514f119588cfa933f354c1"
		break

	case exchange.IBANKDIGITAL:
		config.API_KEY = "9a5d3da7-b1rkuf4drg-dd1aaf13-a8ecb"
		config.API_SECRET = "07d441ab-62e8d245-45d78833-6ca1c"
		break

	case exchange.LBANK:
		config.API_KEY = "c174776e-9605-4b87-ba25-ef290a1e6897"
		config.API_SECRET = "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAIxbnL67Zt6wOqSmFKnDecHSsk0k8Rhojwnzbn7LV+sN4V7ltO14gVEow3kJDxh3EouWPtvH57Hjx1Ft6d41MRfZGDKSVaJJ0gWZEK5pYvZVY52FkpWGe5uY3PwxujC2yaaXSkTy+FazOZKsxHWrTyuh3MgPQgl3xcO4HqNB65DFAgMBAAECgYBc0CPJeFDhBvXwdKaLT+jOw44GN1x6gIG92cyCaeKcW5RhVVKcCaixy1vfSJ9D1VFdHqA4Y2uSFYZzEVSqDNCF7wS1OeiBmec0CcHGs40HVJHAZJFwVfBOLUdio33LYqoUDNzqtbEGmkGs1RdQLj7o9JhFUu4sEN4LfO05JWprIQJBAOH3B5BuWsrPcoLuBJVeePpkFdkJ65zVboDakhDziluieigqqzT+kY+U0ckn4NuOs/ueKa5CwXcTgSqIXCRg2lkCQQCfA5th2YGw/fGUAcgADST38KHb0OW+kfF+H63lZPISXiXIbbbBL1EpYjtn2WTsrcadH2BB9rpIdMyoUZszVIRNAkBrRWKJ5lmjvieWkHgMkPTNqYXVqyf3JDt5YEnHUlZ0egWT2+27Er73cqbE3/GXSX+YC9WtrHM7nD7Nej6D5pbBAkBE8G5kLMWCc4ZR0bfg9dHqQIQb5eRFC8b0FE3zHyGn/vNIgvBxrs70LydsLZ8I0YpDQoAb+RjoIuM7si2kQmcdAkEAsG2ywEw+2wrswodWFDjbGhVVQH8Oow1lPkM2ljj74e4da0ngSLkDKl7taSL64c7l/jyGacR+pXk5NC8IbmegXA=="
		break

	case exchange.BITMART:
		config.API_KEY = ""
		config.API_SECRET = ""
		config.Passphrase = "" // key name
		break

	case exchange.BIKI:
		config.API_KEY = "0b4d7ee72e98a137aa6d2feb7d29fb4b"    //"e93427e0fb84b8e76ecefb78a560cf11"
		config.API_SECRET = "819061f36cc4b7f55c7f398629d1201e" //"82b3e215e921062987aae836c87fe3a0"
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
		config.API_KEY = ""
		config.API_SECRET = ""
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
		config.API_KEY = ""
		config.API_SECRET = ""
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
		config.API_KEY = "cgdvfeka"
		config.API_SECRET = "rns93eby"
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

	case exchange.BITHUMB:
		config.API_KEY = "73e2cd0b89a63ab4ccf67fcf00f5bb96"
		config.API_SECRET = "6b1b05bb851d7e7d2ddc26428757aaf250540821a2c93d333cb759b314dd403c"

	case exchange.SWITCHEO:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.BLOCKTRADE:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.BKEX:
		config.API_KEY = "9f961a84171d245ad86b3ee9b88dc3898686fa0be9320a6f228b1589cac9ed04"    //"c439e89707e2fdf73e61c3dd029aede2379c2b598fba8da07967937e30038344"
		config.API_SECRET = "769651b25bb0ae728e8e0675ca4f1b6ac5fdafc1a37e1e27ee80ff52e439f191" //"a3cbd244c25c286ca5a2aa45255a97bb29840bdd9a331788f53278ad2a7742fa"

	case exchange.NEWCAPITAL:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.COINDEAL:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.HIBITEX:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.BGOGO:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.FTX:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.TXBIT:
		config.API_KEY = "f3d0dbf63000BD3DfcbCd46bDc4C8f2a"
		config.API_SECRET = "0DB1bA1Fb2aed98b8C2b7eF8Bf2a6E0f"

	case exchange.PROBIT:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.BITPIE:
		config.API_KEY = ""
		config.API_SECRET = ""

	case exchange.TAGZ:
		config.API_KEY = "1"
		config.API_SECRET = "3"

	case exchange.IDCM:
		config.API_KEY = "sYpHrC1vl0uK2vW0aACunw"
		config.API_SECRET = "Xhw6V3pCYFSWsOugbcFA2u6nLi30IK8W1bbVXRkEhTKvHWY6JEw6BFcsb8lCZoTexVUhCzijS0i8JBFOkXh4d1kmZJKOwHESZwYL2HtalceAm4IAtwi27P0wAyiTOuNamq8dCTMDu7I1NKZU09l0pCLavratpzlRo3eWVvXzJEy07ZIQznTF9NLJEYIbJjeULQCJv2ZucLgK4K2vJjU1FGIBivhXNMQqqDehmYlwopsnw0EUVcrpoce581r0K1Hq"

	case exchange.HOO:
		config.API_KEY = "3aHUw1fBXu56xCYTLYfdoapoLXpeiy"
		config.API_SECRET = "Kwmsx4XgabpfTfvR6PQ1b9LcHPHUoDt4UDyA5TzNEXgpMeVWed"

	case exchange.HOMIEX:
		config.API_KEY = "vaqHzu0GyduaYy8QgMU6Y2JA98revIDUNl7b2sVYPUeOch5vSvX5mNPWOb1K7pgC"
		config.API_SECRET = "NqrSXeHc1KWPFhD1iQYiINCHL2xbWfyqBNb0m1Se5nekhFJv0LXhce1kk2drReLp"
	}
}
