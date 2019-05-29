# Contributing to GoRED

:+1: First off, thanks for taking the time to contribute! :+1:

The following is a set of guidelines for contributing to GoRED and its packages, which are hosted in the [Bitontop Technologies Inc.](https://github.com/bitontop) on GitHub. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

** Duplicate Folder **
* Go to the path `~/exchange`


Data Collection

1.0 New Exchange Platform

1.1 Develop the Basic Functions

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

    1.1.5 Add Exchange Information
        1.1.5.1 Add exchange information in file: 
            gored/test/conf/apikey.go
            gored/main.go
            gored/initial/iniman.go
            gored/exchange/meta.go  ---→  find a new ID and set to DEFAULT_ID in file [gored/exchange/meta.go]

//-----------↓ todo ------------

1.2 Get RealTime Data

    1.2.1 Init Exchange Config
        1.2.1.1 Modify [main.go]
            1.2.1.1.1 Add Function [init"ExchangeName"()]
            1.2.1.1.2 Modify Config Content
            1.2.1.1.3 Add [init"ExchangeName"()] in Init()
            
    1.2.2 Initial New Exchange Task
        1.2.2.1 Modify [init_task.go]
            1.2.2.1.1 Get Exchange Config [e"ExchangeName" := exMan.Get(exchange."EXCHANGENAME")]
            1.2.2.1.2 Call InitTask Function [m.InitTask(e"ExchangeName".GetPairs(), exchange."EXCHANGENAME", [pairs_amount])]
            
    1.2.3 Implement Gaining RealTime Data
        1.2.3.1 Modify [data.go]
            1.2.3.1.1 Get Exchange Config [e"ExchangeName" := exMan.Get(exchange."EXCHANGENAME")]
            1.2.3.1.2 Call InitTask Function [m.InitTask(e"ExchangeName".GetPairs(), exchange."EXCHANGENAME", [pairs_amount])]
            
    1.2.4 Deploy the program on Server