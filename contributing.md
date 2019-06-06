# Contributing to GoRED

:+1: First off, thanks for taking the time to contribute! :+1:

The following is a set of guidelines for contributing to GoRED and its packages, which are hosted in the [Bitontop Technologies Inc.](https://github.com/bitontop) on GitHub. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

** Duplicate Folder **
* Go to the path `~/exchange`


Data Collection

1.0 New Exchange Platform

1.1 Develop the Basic Functions
    1.1.0 prepare
        1.1.0.1 fork gored to your repo
        1.1.0.2 in VScode terminal, go get -u "github.com/bitontop/gored"
        1.1.0.3 using git pull to update gored code periodically
    1.1.1 Duplicate Exchange Template
        1.1.1.1 Duplicate Folder [/gored/exchange/blank]
        1.1.1.2 Rename ["blank"] folder -> ["exchange name"] and Put the folder under [/gored/exchange/]
        
    1.1.2 Duplicate Exchange Test Case Template
        1.1.2.1 Duplicate File [/gored/test/blank_test.go]
        1.1.2.2 Rename ["blank"_test.go] -> ["exchange name"_test.go] and Put the file under [/gored/test/]
        1.1.2.3 Replace [blank] to ["exchange name"]
        1.1.2.4 Replace [Blank] to ["Exchange Name"]
        1.1.2.5 Replace [BLANK] to ["EXCHANGE NAME"]
        
    1.1.3 Develop Basic Functions
        1.1.3.1 Follow the instruction on each file which is under [/gored/exchange/"exchange name"]
        1.1.3.2 Complete all public functions 
        1.1.3.3 Complete all private functions 
        
    1.1.4 Test Basic Functions
        1.1.4.1 Run Each Test Case to Make Sure the function is working
        1.1.4.2 Deposit asset for private function test:
            login to account, go to asset
            go into deposit, using the deposit address to deposit
            wait until the conformation complete

2.1 after finish
    2.1.1 add new exchange information/function in file:
        2.1.1.1gored/test/conf/apikey.go:
            in function Exchange, add new case exchange.["EXCHANGE NAME"]
        2.1.1.2gored/main.go:
            add import "github.com/bitontop/gored/exchange/["exchange name"]"
            in function Init, add: Init["Exchange Name"](config)
            add function: func Init["Exchange Name"]
        2.1.1.3gored/initial/iniman.go
            add import "github.com/bitontop/gored/exchange/["exchange name"]"
            in function Init, add new case exchange.["EXCHANGE NAME"]
        2.1.1.4gored/exchange/meta.go
            in function initExchangeNames, add supportList = append(supportList, ["EXCHANGE NAME"]), comment a new ID beside it
            set DEFAULT_ID in file [gored/exchange/meta.go] to this ID


3.1 make pull request to branch: pullRequest

