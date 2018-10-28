function Wallet() {

    var componentCounter = 0;

    var seed = null;

    var Waves = null;

    // Payment method
    this.pay = function() {
        var addressTo = getEl('addressTo').value;
        var amount = getEl('amount').value;
        var currency = getEl('paymentCurrency').selectedIndex;
        if (validatePaymentFields(addressTo, amount)) {
            switch (currency) {
                case 2:
                    if (checkAddress(addressTo)) {
                        $('#outPayMsgBtc').show();
                        $('#outPayMsgEth').hide();
                    } else {
                        $('#outPayMsgBtc').hide();
                        $('#outPayMsgEth').hide();
                    }
                    break;
                case 3:
                    var web3 = new Web3(Web3.currentProvider);
                    if (web3.isAddress(addressTo)) {
                        $('#outPayMsgEth').show();
                        $('#outPayMsgBtc').hide();
                    } else {
                        $('#outPayMsgEth').hide();
                        $('#outPayMsgBtc').hide();
                    }
                    break;
                default:
                    $('#outPayMsgEth').hide();
                    $('#outPayMsgBtc').hide();
            }
            $("#modalPay").modal();
        }
    }

    // Payment method
    this.payConfirmed = function() {
        $("#modalPay").modal("hide");
        var addressTo = getEl('addressTo').value;
        var amount = new Decimal(getEl('amount').value);
        var currency = getEl('paymentCurrency').selectedIndex;
        var feeCurrency = getEl('payFeeCurrency').selectedIndex;

        if (validatePaymentFields(addressTo, amount)) {
            switch (currency) {
                case 1:
                    transfer(addressTo, amount, 'WAVES', '', feeCurrency);
                    break;
                case 2:
                    if (checkAddress(addressTo)) {
                        var fee = new Decimal(0.0005);
                        transfer(nodeAddress, amount.add(fee).toNumber(), '7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs', 'forwardbtc=' + addressTo, feeCurrency);
                    } else {
                        transfer(addressTo, amount, '7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs', '', feeCurrency);
                    }
                    break;
                case 3:
                    var web3 = new Web3(Web3.currentProvider);
                    if (web3.isAddress(addressTo)) {
                        var fee = new Decimal(0.001);
                        transfer(nodeAddress, amount.add(fee).toNumber(), '4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53', 'forwardeth=' + addressTo, feeCurrency);
                    } else {
                        transfer(addressTo, amount, '4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53', '', feeCurrency);
                    }
                    break;
                default:
                    transfer(addressTo, amount, '4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf', '', feeCurrency);
            }

            getEl('addressTo').value = '';
            getEl('amount').value = '';

            $('#payAdvanced').fadeOut();
        }
    }

    // Copy method for copying address to clipboard
    this.copy = function(elId) {
        var copyText = document.getElementById(elId);
        copyText.select();
        document.execCommand("copy");
        $('#copymessage' + elId).fadeIn(function() {
            setTimeout(() => {
                $('#copymessage' + elId).fadeOut();
            }, 2000);
        });
    }

    // Saves user's nickname to blockchain
    this.save = function() {
        var nickname = getEl('nickname').value;
        var email = getEl('email').value;
        var country = getEl('countryHidden').value;
        var city = getEl('cityHidden').value;
        if (validateSettingsFields(nickname, email, country, city)) {
            $.ajax({
                url: '/apply/',
                method: 'POST',
                data: {
                    nickname: nickname,
                    email: email,
                    country: country,
                    city: city
                },
                success: function (data, status){
                    $('#applicationForm').fadeOut(function (){
                        $('#applicationSuccess').fadeIn();
                    });
                },
                error: function(data, status, error) {
                    if (data.responseJSON.message.search('nickname') != -1) {
                        setHTML('errorMessageNickname', allLocales.nicknameExists);
                        $('#nicknameGroup').addClass('has-error');
                        $('#errorMessageNickname').fadeIn(function() {
                            setTimeout(() => {
                                $('#errorMessageNickname').fadeOut();
                                $('#nicknameGroup').removeClass('has-error');
                            }, 2000);
                        });
                    } else if (data.responseJSON.message.search('email') != - 1) {
                        setHTML('errorMessageNickname', allLocales.emailExists);
                        $('#emailGroup').addClass('has-error');
                        $('#errorMessageNickname').fadeIn(function() {
                            setTimeout(() => {
                                $('#errorMessageNickname').fadeOut();
                                $('#emailGroup').removeClass('has-error');
                            }, 2000);
                        });
                    } else {
                        setHTML('errorMessageNickname', data.responseJSON.message);
                        $('#errorMessageNickname').fadeIn(function() {
                            setTimeout(() => {
                                $('#errorMessageNickname').fadeOut();
                            }, 2000);
                        });
                    }
                } 
            });
        }
    }

    // Updates hidden country and city fields
    this.updateCountryCity = function() {
        var country = $('#country').val();
        var city = $('#city').val();

        if (country == 'Country') {
            country = '';
        }

        if (city == 'City') {
            city = '';
        }

        $('#countryHidden').val(country);
        $('#cityHidden').val(city);
    }

    // Exchange currencies
    this.exchange = function() {
        var selectedCurrency = getEl('currency').selectedIndex;
        var amount = getEl('amount').value;
        if (validateExchangeFields(selectedCurrency, amount)) {
            $("#modalExchange").modal();
        }
    }

    // Confirms exchange currencies
    this.exchangeConfirm = function() {
        $("#modalExchange").modal('hide');
        var selectedCurrency = getEl('currency').selectedIndex;
        var amount = getEl('amount').value;
        if (validateExchangeFields(selectedCurrency, amount)) {
            $('#content').fadeOut(function() {
                $('#transactionInProgress').fadeIn();
            });
            switch (selectedCurrency) {
                case 1:
                    btcToAno(amount);
                    break;
                case 2:
                    ethToAno(amount);
                    break;
                default:
                    wavToAno(amount);
            }
        }
    }

    // Withdraw ANT profit
    this.withdraw = function() {
        $("#modalWithdraw").modal();
    }

    // Withdraw confirmed
    this.withdrawConfirm = function() {
        $("#modalWithdraw").modal('hide');
        var withdrawCurrency = getEl('withdrawFeeCurrency').selectedIndex;
        if (withdrawCurrency == 0) {
            var assetId = '4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf';
            var amount = 0.9;
        } else if (feeCurrency == 2) {
            var assetId = '7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs';
            var amount = 0.0000012;
        } else if (feeCurrency == 3) {
            var assetId = '4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53';
            var amount = 0.000024;
        } else {
            var assetId = 'WAVES';
            var amount = 0.003;
        }
        transfer(nodeAddress, amount, assetId, 'withdraw', withdrawCurrency);
    }

    // Sign in method
    this.signIn = function() {
        var pass = getEl('password').value;
        if (validatePasswordField(pass)) {
            try {
                var restoredPhrase = Waves.Seed.decryptSeedPhrase(window.localStorage.getItem('encrypted'), pass);
                window.localStorage.setItem('seed', restoredPhrase);
                seed = Waves.Seed.fromExistingPhrase(restoredPhrase);
                Cookies.set('address', seed.address, { expires: 365 });
                $('#signInForm').fadeOut(function() {
                    $('#newGroupAjax').fadeIn();
                });
                $.ajax({
                    url: '/sign-in/',
                    method: 'POST',
                    data: {
                        password: pass
                    },
                    success: function(data, status) {
                        window.location = "/";
                    },
                    error: function(data, status, error) {
                        console.log(error);
                    }
                });
            } catch (e) {
                setHTML('required', e);
                $('#required').fadeIn(function() {
                    setTimeout(() => {
                        $('#required').fadeOut();
                        $('#passwordGroup').removeClass('has-error');
                    }, 2000);
                });
            }
        }
    }

    // Sign out method
    this.signOut = function() {
        window.localStorage.removeItem('seed');
        Cookies.remove('address');
        window.location = '/sign-out/';
    }

    // Sign up next method
    this.signUpNext = function() {
        setValue('seedinput', allLocales.copyAgain);
        var copyText = document.getElementById("seedinput");
        copyText.select();
        document.execCommand("copy");
        $('#newGroup').fadeOut(function() {
            $('#newGroup1').fadeIn();
        });
    }

    // Sign up next method (second step)
    this.signUpNext1 = function() {
        var seedTa = getEl('seedta').value;
        if (validateSeedField(seedTa, true)) {
            $('#newGroup1').fadeOut(function() {
                $('#newGroup2').fadeIn();
            });
        }
    }

    // Sign up next method (third step)
    this.signUpNext2 = function() {
        var p1 = getEl('password1').value;
        var p2 = getEl('password2').value;
        if (validateSUPasswords(p1, p2)) {
            window.localStorage.setItem('seed', seed.phrase);
            window.localStorage.setItem('encrypted', seed.encrypt(p1));
            Cookies.set('address', seed.address, { expires: 365 });
            $('#newGroup2').fadeOut(function() {
                $('#newGroupAjax').fadeIn();
            });
            $.ajax({
                url: '/sign-up/',
                method: 'POST',
                data : {
                    password: p1
                },
                success: function(data, status) {
                    window.location = "/";
                },
                error: function(data, status, error) {
                    console.log(error);
                }
            });
        }
    }

    // Sign out copy method
    this.signUpCopy = function() {
        var copyText = document.getElementById("seedinput");
        copyText.select();
        document.execCommand("copy");
        $('#copymessage').fadeIn(function() {
            setTimeout(() => {
                $('#copymessage').fadeOut();
            }, 2000);
        });
    }

    // Import next method
    this.importNext = function() {
        var seedTa = getEl('seedta').value;
        if (validateSeedField(seedTa, false)) {
            seed = Waves.Seed.fromExistingPhrase(seedTa);
            $('#importGroup1').fadeOut(function() {
                $('#importGroup2').fadeIn();
            });
        }
    }

    // Import next method (second step)
    this.importNext1 = function() {
        var p1 = getEl('password1').value;
        var p2 = getEl('password2').value;
        if (validateSUPasswords(p1, p2)) {
            window.localStorage.setItem('seed', seed.phrase);
            window.localStorage.setItem('encrypted', seed.encrypt(p1));
            Cookies.set('address', seed.address, { expires: 365 });
            $('#importGroup2').fadeOut(function() {
                $('#newGroupAjax').fadeIn();
            });
            $.ajax({
                url: '/sign-up/',
                method: 'POST',
                data : {
                    password: p1
                },
                success: function(data, status) {
                    window.location = "/";
                },
                error: function(data, status, error) {
                    console.log(error);
                }
            });
        }
    }

    this.updateCalculator = function() {
        var amount = getEl('amount').value;
        var sel = getEl('currency').selectedIndex;
        switch(sel) {
            case 1:
                var price = priceBtc;
                break;
            case 2:
                var price = priceEth;
                break;
            default:
                var price = priceWav;
        }

        var totalAmountAno = 0;
        var priceFactor = price * (anotePriceFactor / anotePrice);

        while (amount > 0) {
            var amountAno = amount / price;
            console.log(amountAno);
            if (amountAno > (anoteTierPrice / 100000000)) {
                amountAno = anoteTierPrice / 100000000;
            }

            amount -= amountAno * price;

            totalAmountAno += amountAno;
            price += priceFactor;
        }

        // var amountAno = (amount / price).toFixed(2);
        totalAmountAno = totalAmountAno.toFixed(2)
        var amountEur = (totalAmountAno * 17).toFixed(2);
        setHTML('amountAno', totalAmountAno);
        setHTML('amountEur', amountEur);

    }

    this.showSeed = function() {
        // alert(seed.phrase);
        getEl('seed').innerHTML = seed.phrase;
        getEl('copySeed').disabled = false;
    }

    // PRIVATE METHODS

    // Constructor method
    function constructor() {
        $(".sidebar-menu a").each(function() {
            if ($(this).attr('href') == window.location.pathname) {
                $(this).parent().addClass('active');
            }
        });

        var referral = getReferralFromUrl();

        if (referral) {
            Cookies.set('referral', referral, { expires: 30, domain: getDomainName(window.location.hostname) });
        }

        Waves = WavesAPI.create(WavesAPI.MAINNET_CONFIG);

        switch(window.location.pathname) {
            case '/sign-in/':
                initSignIn();
                break;
            case '/sign-up-new/':
                newWallet();
                break;
            case '/':
                initSuccess();
                break;
        }

        var encrypted = window.localStorage.getItem('encrypted');
        var restoredPhrase = window.localStorage.getItem('seed');

        if (restoredPhrase) {
            seed = Waves.Seed.fromExistingPhrase(restoredPhrase);
        } else if (encrypted && !window.location.href.endsWith('/sign-in/')) {
            window.location.href = '/sign-in/';
        } else if (!encrypted && !window.location.href.endsWith('/sign-up/') && !window.location.href.endsWith('/sign-up-new/') && !window.location.href.endsWith('/sign-up-import/')) {
            window.location.href = '/sign-up/';
        }

        if (getEl('transactionsButton')) {
            getEl('transactionsButton').href += address;
        }
    }

    // Successful init
    function initSuccess() {
        setValue('address', address);

        Waves.API.Node.v1.assets.balance(address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf").then((balance) => {
            var antBalance = parseFloat(parseFloat(balance.balance) / parseFloat(10**8)).toFixed(5);
            setHTML('balanceAnt', antBalance);
            updateCounter();
        });

        Waves.API.Node.v1.assets.balance(address, "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs").then((balance) => {
            var btcBalance = parseFloat(parseFloat(balance.balance) / parseFloat(10**8)).toFixed(5);
            setHTML('balanceBtc', btcBalance);
            updateCounter();
        });

        Waves.API.Node.v1.assets.balance(address, "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53").then((balance) => {
            var ethBalance = parseFloat(parseFloat(balance.balance) / parseFloat(10**8)).toFixed(5);
            setHTML('balanceEth', ethBalance);
            updateCounter();
        });

        Waves.API.Node.v1.addresses.balance(address).then((balance) => {
            var wavBalance = parseFloat(parseFloat(balance.balance) / parseFloat(10**8)).toFixed(5);
            setHTML('balanceWav', wavBalance);
            updateCounter();
        });

        timeout = setTimeout(initSuccess, 1000);
    }

    // Init sign in page
    function initSignIn() {
        var encrypted = window.localStorage.getItem('encrypted');
        if (!encrypted) {
            window.location.href = '/sign-up/';
        }
    }

    // Updates counter for loading purposes
    function updateCounter() {
        componentCounter++;
        if (componentCounter == 4) {
            new QRious({
                size: 300,
                element: document.getElementById('qr'),
                value: 'waves://' + address
            });

            new QRious({
                size: 300,
                element: document.getElementById('qrbitcoin'),
                value: 'bitcoin://' + getEl('addressBitcoin').value
            });

            new QRious({
                size: 300,
                element: document.getElementById('qrethereum'),
                value: getEl('addressEthereum').value
            });

            $('#loader').fadeOut(function() {
                $('#content').fadeIn();
            });
        }
    }

    // Gets DOM element
    function getEl(id) {
        return document.getElementById(id);
    }

    // Sets HTML for element
    function setHTML(id, html) {
        getEl(id).innerHTML = html;
    }

    // Sets value for element
    function setValue(id, value) {
        getEl(id).value = value;
    }

    // Binds event the right way
    function bind(scope, fn) {
        return function() {
            return fn.apply(scope, arguments);
        }
    }

    // Checks and validates fields for payments
    function validatePaymentFields(addressTo, amount) {
        var validates = true;

        if (addressTo.length == 0) {
            $('#addressToGroup').addClass('has-error');
            validates = false;
        }

        if (amount.length == 0) {
            $('#amountGroup').addClass('has-error');
            validates = false;
        }

        if (!validates) {
            setHTML('errorMessagePayment', allLocales.bothFields);
            $('#errorMessagePayment').fadeIn(function() {
            setTimeout(() => {
                $('#errorMessagePayment').fadeOut();
                $('#amountGroup').removeClass('has-error');
                $('#addressToGroup').removeClass('has-error');
            }, 2000);
        });
        }

        return validates;
    }

    // Checks and validates fields for settings form
    function validateSettingsFields(nickname, email, country, city) {
        var validates = true;

        if (nickname.length == 0) {
            $('#nicknameGroup').addClass('has-error');
            validates = false;
        }

        if (email.length == 0) {
            $('#emailGroup').addClass('has-error');
            validates = false;
        }

        if (country.length == 0) {
            $('#countryGroup').addClass('has-error');
            validates = false;
        }

        if (city.length == 0) {
            $('#cityGroup').addClass('has-error');
            validates = false;
        }

        if (!validates) {
            setHTML('errorMessageNickname', allLocales.allFields);
            $('#errorMessageNickname').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessageNickname').fadeOut();
                    $('#nicknameGroup').removeClass('has-error');
                    $('#emailGroup').removeClass('has-error');
                    $('#countryGroup').removeClass('has-error');
                    $('#cityGroup').removeClass('has-error');
                }, 2000);
            });
        }

        return validates;
    }

    // Checks and validates field for password form
    function validatePasswordField(password) {
        var validates = true;

        if (password.length == 0) {
            $('#passwordGroup').addClass('has-error');
            validates = false;
        }

        if (!validates) {
            setHTML('required', allLocales.jsFieldEmpty);
            $('#required').fadeIn(function() {
                setTimeout(() => {
                    $('#required').fadeOut();
                    $('#passwordGroup').removeClass('has-error');
                }, 2000);
            });
        }

        return validates;
    }

    // Checks and validates fields for exchange form
    function validateExchangeFields(selectedCurrency, amount) {
        var validates = true;

        if (amount.length == 0) {
            $('#amountGroup').addClass('has-error');
            validates = false;
        }

        if (!validates) {
            setHTML('errorMessageExchange', allLocales.jsFieldEmpty);
            $('#errorMessageExchange').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessageExchange').fadeOut();
                    $('#amountGroup').removeClass('has-error');
                }, 2000);
            });
        }

        return validates;
    }

    // Checks and validates fields for seed
    function validateSeedField(seedTa, checkSeed) {
        var validates = true;

        if (seedTa.length == 0) {
            $('#seedGroup').addClass('has-error');
            validates = false;
            setHTML('errorMessageSeed', allLocales.jsFieldEmpty);
            $('#errorMessageSeed').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessageSeed').fadeOut();
                    $('#seedGroup').removeClass('has-error');
                }, 2000);
            });
        }

        if (validates && checkSeed && seedTa != seed.phrase) {
            $('#seedGroup').addClass('has-error');
            validates = false;
            setHTML('errorMessageSeed', allLocales.wrongSeed);
            $('#errorMessageSeed').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessageSeed').fadeOut();
                    $('#seedGroup').removeClass('has-error');
                }, 2000);
            });
        }

        return validates;
    }

    // Checks and validates sign up password fields
    function validateSUPasswords(p1, p2) {
        var validates = true;

        if (p1.length == 0) {
            $('#passwordGroup1').addClass('has-error');
            validates = false;
            setHTML('errorMessagePassword', allLocales.bothFields);
            $('#errorMessagePassword').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessagePassword').fadeOut();
                    $('#passwordGroup1').removeClass('has-error');
                }, 2000);
            });
        }

        if (p2.length == 0) {
            $('#passwordGroup2').addClass('has-error');
            validates = false;
            setHTML('errorMessagePassword', allLocales.bothFields);
            $('#errorMessagePassword').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessagePassword').fadeOut();
                    $('#passwordGroup2').removeClass('has-error');
                }, 2000);
            });
        }

        if (validates && p1 != p2) {
            $('#passwordGroup2').addClass('has-error');
            validates = false;
            setHTML('errorMessagePassword', allLocales.passNotMatch);
            $('#errorMessagePassword').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessagePassword').fadeOut();
                    $('#passwordGroup2').removeClass('has-error');
                }, 2000);
            });
        }

        return validates;
    }

    // Transfers any token
    function transfer(addressTo, amount, assetId, attachment, feeCurrency) {
        if (feeCurrency == 1) {
            var feeAssetId = 'WAVES';
            var feeAmount = 100000;
        } else if (feeCurrency == 2) {
            var feeAssetId = '7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs';
            var feeAmount = 40;
        } else if (feeCurrency == 3) {
            var feeAssetId = '4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53';
            var feeAmount = 800;
        } else {
            var feeAssetId = '4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf';
            var feeAmount = 30000000;
        }

        $('#content').fadeOut(function() {
            $('#transactionInProgress').fadeIn(function() {
                const transferData = {
                    recipient: addressTo,
                    assetId: assetId,
                    amount: parseInt(amount * 10**8),
                    feeAssetId: feeAssetId,
                    fee: feeAmount,
                    attachment: attachment,
                    timestamp: Date.now()
                };

                Waves.API.Node.v1.assets.transfer(transferData, seed.keyPair).then((responseData) => {
                    $('#transactionInProgress').fadeOut(function() {
                        $('#transactionSuccess').fadeIn(function() {
                            setTimeout(() => {
                                $('#transactionSuccess').fadeOut();
                            }, 10000);
                        });
                        $('#content').fadeIn();
                    });
                }).catch((err) => {
                    setHTML('errorMessage', err);
                    $('#transactionInProgress').fadeOut(function() {
                        $('#transactionError').fadeIn(function() {
                            setTimeout(() => {
                                $('#transactionError').fadeOut();
                            }, 10000);
                        });
                        $('#content').fadeIn();
                    });
                });
            });
        });
    }

    // Exchange WAV to ANO
    function wavToAno( amount) {
        var feeCurrency = getEl('exchangeFeeCurrency').selectedIndex;
        transfer(nodeAddress, amount, 'WAVES', '', feeCurrency);
    }

    // Exchange BTC to ANO
    function btcToAno(amount) {
        var feeCurrency = getEl('exchangeFeeCurrency').selectedIndex;
        transfer(nodeAddress, amount, '7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs', '', feeCurrency);
    }

    // Exchange ETH to ANO
    function ethToAno(amount) {
        var feeCurrency = getEl('exchangeFeeCurrency').selectedIndex;
        transfer(nodeAddress, amount, '4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53', '', feeCurrency);
    }

    // Gets referral from url
    function getReferralFromUrl(){
        var k = 'r';
        var p={};
        location.search.replace(/[?&]+([^=&]+)=([^&]*)/gi,function(s,k,v){p[k]=v})
        return k?p[k]:p;
    }

    // Gets domain name from url
    function getDomainName(hostName) {
        return hostName.substring(hostName.lastIndexOf(".", hostName.lastIndexOf(".") - 1) + 1);
    }

    // New wallet method
    function newWallet() {
        seed = Waves.Seed.create();
        setHTML('seed', seed.phrase);
        setValue('seedinput', seed.phrase);
        $('#newGroup').fadeIn();
    }

    // Attach all events
    switch (window.location.pathname) {
        case '/':
            getEl('payButton').addEventListener('click', bind(this, this.pay), false);
            getEl('payConfirmButton').addEventListener('click', bind(this, this.payConfirmed), false);
            break;
        case '/profit/':
            getEl('withdrawButton').addEventListener('click', bind(this, this.withdraw), false);
            getEl('withdrawConfirmButton').addEventListener('click', bind(this, this.withdrawConfirm), false);
            break;
        case '/settings/':
            if (getEl('saveButton')) {
                getEl('saveButton').addEventListener('click', bind(this, this.save), false);
                getEl('country').addEventListener('change', bind(this, this.updateCountryCity), false);
                getEl('city').addEventListener('change', bind(this, this.updateCountryCity), false);
            }
            getEl('showSeed').addEventListener('click', bind(this, this.showSeed), false);
            break;
        case '/exchange/':
            getEl('exchangeButton').addEventListener('click', bind(this, this.exchange), false);
            getEl('exchangeConfirmButton').addEventListener('click', bind(this, this.exchangeConfirm), false);
            getEl('currency').addEventListener('change', bind(this, this.updateCalculator), false);
            getEl('amount').addEventListener('keyup', bind(this, this.updateCalculator), false);
            break;
        case '/sign-in/':
            getEl('signInButton').addEventListener('click', bind(this, this.signIn), false);
            break;
        case '/sign-up/':
            break;
        case '/sign-up-new/':
            getEl('signupnext').addEventListener('click', bind(this, this.signUpNext), false);
            getEl('signupnext1').addEventListener('click', bind(this, this.signUpNext1), false);
            getEl('signupcopy').addEventListener('click', bind(this, this.signUpCopy), false);
            break;
        case '/sign-up-import/':
            getEl('importnext').addEventListener('click', bind(this, this.importNext), false);
            getEl('importnext1').addEventListener('click', bind(this, this.importNext1), false);
            break;
    }

    // Calling Wallet constructor
    constructor();
}
