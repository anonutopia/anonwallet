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
                        setHTML('errorMessageSettings', allLocales.nicknameExists);
                        $('#nicknameGroup').addClass('has-error');
                        $('#errorMessageSettings').fadeIn(function() {
                            setTimeout(() => {
                                $('#errorMessageSettings').fadeOut();
                                $('#nicknameGroup').removeClass('has-error');
                            }, 2000);
                        });
                    } else if (data.responseJSON.message.search('email') != - 1) {
                        setHTML('errorMessageSettings', allLocales.emailExists);
                        $('#emailGroup').addClass('has-error');
                        $('#errorMessageSettings').fadeIn(function() {
                            setTimeout(() => {
                                $('#errorMessageSettings').fadeOut();
                                $('#emailGroup').removeClass('has-error');
                            }, 2000);
                        });
                    } else {
                        setHTML('errorMessageSettings', data.responseJSON.message);
                        $('#errorMessageSettings').fadeIn(function() {
                            setTimeout(() => {
                                $('#errorMessageSettings').fadeOut();
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

    // Sign out method
    this.signOut = function() {
        window.sessionStorage.removeItem('seed');
        window.location = '/sign-out/';
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
        getEl('seed').innerHTML = seed.phrase;
        getEl('copySeed').disabled = false;
    }

    this.signUpSeedChange = function() {
        var seedPhrase = $('#seed').val();
        try {
            seed = Waves.Seed.fromExistingPhrase(seedPhrase);
            setValue('address', seed.address);
        } catch (e) {}
    }

    this.checkPassword = function() {
        $('#seed').val('');
        $('#messageError').hide();
        $('#messageSuccess').hide();

        var password = $('#passwordold').val();
        var encrypted = window.localStorage.getItem('encrypted');

        try {
            var seedPhrase = Waves.Seed.decryptSeedPhrase(encrypted, password);
            seed = Waves.Seed.fromExistingPhrase(seedPhrase);
            if (seed.address == address) {
                $('#messageError').hide();
                $('#messageSuccess').html("Success: The password is right.");
                $('#messageSuccess').show();
                $('#address').val(seed.address);
                $('#seed').val(seed.phrase);
            } else {
                throw new Error('The address on this seed is not right.');
            }
        } catch (e) {
            $('#messageSuccess').hide();
            $('#messageError').html(e + '.');
            $('#messageError').show();
        }
    }

    this.checkSeed = function() {
        $('#passwordold').val('');
        $('#messageError').hide();
        $('#messageSuccess').hide();

        var seedPhrase = $('#seed').val();

        try {
            seed = Waves.Seed.fromExistingPhrase(seedPhrase);
            console.log(address);
            console.log(seed.address);
            if (seed.address == address) {
                $('#messageError').hide();
                $('#messageSuccess').html("Success: The seed is right.");
                $('#messageSuccess').show();
                $('#address').val(seed.address);
            } else {
                throw new Error('The address on this seed is not right.');
            }
        } catch (e) {
            $('#messageSuccess').hide();
            $('#messageError').html(e + '.');
            $('#messageError').show();
        }
    }

    this.passwordReset = function() {
        var password = $('#password').val();
        var addr = $('#address').val();
        if (password.length > 0) {
            if (addr == address) {
                var seedPhrase = $('#seed').val();

                if (seedPhrase.length > 0) {
                    seed = Waves.Seed.fromExistingPhrase(seedPhrase);
                    window.localStorage.setItem('encrypted', seed.encrypt(password));
                    window.localStorage.setItem('address', address);
                    window.sessionStorage.setItem('seed', seedPhrase);

                    return true;
                }
            } else {
                $('#messageSuccess').hide();
                $('#messageError').html('Something is wrong with your seed or password.');
                $('#messageError').show();
            }
        } else {
            $('#messageSuccess').hide();
            $('#messageError').html("New password can't be empty.");
            $('#messageError').show();
        }

        return false;
    }

    // PRIVATE METHODS

    // Constructor method
    function constructor() {
        if (address && address.length) {
            window.localStorage.setItem('address', address);
        }

        if (email && email.length) {
            window.localStorage.setItem('email', email);
        }

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
            case '/sign-up/':
                initNewWallet();
                break;
            case '/sign-in/':
                initSignIn();
                break;
            case '/sign-in-clean/':
                initSignInClean();
                break;
            case '/password-reset-finish/':
                initPasswordResetFinish();
                break;
            case '/':
                initSuccess();
                break;
        }

        if (getEl('transactionsButton')) {
            getEl('transactionsButton').href += address;
        }

        var seedPhrase = window.sessionStorage.getItem('seed');
        if (seedPhrase && seedPhrase.length) {
            seed = Waves.Seed.fromExistingPhrase(seedPhrase);
        } else if (window.location.pathname == '/' || window.location.pathname == '/profit/' || window.location.pathname == '/exchange/' || window.location.pathname == '/settings/') {
            window.location.href = '/sign-out/'
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

        // this is for legacy users
        var seedPhrase = window.localStorage.getItem('seed');
        if (seedPhrase && seedPhrase.length > 0) {
            // window.localStorage.removeItem('seed');
            var tempseed = Waves.Seed.fromExistingPhrase(seedPhrase);
            window.localStorage.setItem('address', tempseed.address);
        }

        if (!encrypted) {
            window.location.href = '/sign-in-clean/';
        }

        if (finished) {
            var password = $('#password').val();
            var seedPhrase = Waves.Seed.decryptSeedPhrase(encrypted, password);
            window.sessionStorage.setItem('seed', seedPhrase);

            setTimeout(function() {
                window.location.href = '/';
            }, 3000);
        } else {
            var address = window.localStorage.getItem('address');
            setValue('address', address);
        }
    }

    function initSignInClean() {
        var seed = $('#seed').val();
        var password = $('#password').val();

        if (finished) {
            window.sessionStorage.setItem('seed', seed);
            seed = Waves.Seed.fromExistingPhrase(seed);
            window.localStorage.setItem('encrypted', seed.encrypt(password));
            window.localStorage.setItem('address', seed.address);

            setTimeout(function() {
                window.location.href = '/';
            }, 3000);
        }
    }

    function initPasswordResetFinish() {
        var seedPhrase = window.sessionStorage.getItem('seed');
        if (!seedPhrase) {
            seedPhrase = window.localStorage.getItem('seed');
        }
        var encrypted = window.localStorage.getItem('encrypted');

        if (seedPhrase && seedPhrase.length > 0) {
            $('#seed').val(seedPhrase);
        } else {
            $('#form2').fadeIn();

            if (encrypted) {
                $('#form1').fadeIn();
                $('#message1').fadeIn();
                $('#message3').fadeIn();
            } else {
                $('#message1').fadeIn();
            }
        }

        if (finished) {
            setTimeout(function() {
                window.location.href = '/';
            }, 3000)
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
        var validatesDomain = true;
    
        if (nickname.length == 0) {
            $('#nicknameGroup').addClass('has-error');
            validates = false;
        }

        if (email.length == 0) {
            $('#emailGroup').addClass('has-error');
            validates = false;
        }

        if (!validateEmailDomain(email)) {
            $('#emailGroup').addClass('has-error');
            validatesDomain = false;
        }

        if (country.length == 0) {
            $('#countryGroup').addClass('has-error');
            validates = false;
        }

        if (city.length == 0) {
            $('#cityGroup').addClass('has-error');
            validates = false;
        }

        if (!validates || !validatesDomain) {
            if (!validates) {
                setHTML('errorMessageSettings', allLocales.allFields);
            } else if (!validatesDomain) {
                setHTML('errorMessageSettings', allLocales.emailDomain);
            }
            $('#errorMessageSettings').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessageSettings').fadeOut();
                    $('#nicknameGroup').removeClass('has-error');
                    $('#emailGroup').removeClass('has-error');
                    $('#countryGroup').removeClass('has-error');
                    $('#cityGroup').removeClass('has-error');
                }, 2000);
            });
        }

        return validates && validatesDomain;
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

    // Validates email domain
    function validateEmailDomain(email) {
        for (var i = 0; i < _emailDomains.length; i++) {
            if (email.endsWith(_emailDomains[i])) {
                return true;
            }
        }
        return false;
    }

    function initNewWallet() {
        var seed = $('#seed').val();
        var password = $('#password').val();
        var address = $('#address').val();
        if (seed.length == 0) {
            seed = Waves.Seed.create();
            password = createRandomString(12);
            $('#seed').val(seed.phrase);
            $('#address').val(seed.address);
            $('#password').val(password);
        } else if (finished) {
            window.sessionStorage.setItem('seed', seed);
            seed = Waves.Seed.fromExistingPhrase(seed);
            window.localStorage.setItem('encrypted', seed.encrypt(password));
            window.localStorage.setItem('address', seed.address);

            setTimeout(function() {
                window.location.href = '/';
            }, 3000);
        }
    }

    function createRandomString(length) {   
        var str = "";
        for ( ; str.length < length; str += Math.random().toString( 36 ).substr( 2 ) );
        return str.substr( 0, length );
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
    }

    // Calling Wallet constructor
    constructor();
}

var _emailDomains = [
  /* Default domains included */
  "aol.com", "att.net", "comcast.net", "facebook.com", "gmail.com", "gmx.com", "googlemail.com",
  "google.com", "hotmail.com", "hotmail.co.uk", "mac.com", "me.com", "mail.com", "msn.com",
  "live.com", "sbcglobal.net", "verizon.net", "yahoo.com", "yahoo.co.uk",

  /* Other global domains */
  "email.com", "fastmail.fm", "games.com" /* AOL */, "gmx.net", "hush.com", "hushmail.com", "icloud.com",
  "iname.com", "inbox.com", "lavabit.com", "love.com" /* AOL */, "outlook.com", "pobox.com", "protonmail.com",
  "rocketmail.com" /* Yahoo */, "safe-mail.net", "wow.com" /* AOL */, "ygm.com" /* AOL */,
  "ymail.com" /* Yahoo */, "zoho.com", "yandex.com",

  /* United States ISP domains */
  "bellsouth.net", "charter.net", "cox.net", "earthlink.net", "juno.com",

  /* British ISP domains */
  "btinternet.com", "virginmedia.com", "blueyonder.co.uk", "freeserve.co.uk", "live.co.uk",
  "ntlworld.com", "o2.co.uk", "orange.net", "sky.com", "talktalk.co.uk", "tiscali.co.uk",
  "virgin.net", "wanadoo.co.uk", "bt.com",

  /* Domains used in Asia */
  "sina.com", "sina.cn", "qq.com", "naver.com", "hanmail.net", "daum.net", "nate.com", "yahoo.co.jp", "yahoo.co.kr", "yahoo.co.id", "yahoo.co.in", "yahoo.com.sg", "yahoo.com.ph", "163.com", "126.com", "aliyun.com", "foxmail.com",

  /* French ISP domains */
  "hotmail.fr", "live.fr", "laposte.net", "yahoo.fr", "wanadoo.fr", "orange.fr", "gmx.fr", "sfr.fr", "neuf.fr", "free.fr",

  /* German ISP domains */
  "gmx.de", "hotmail.de", "live.de", "online.de", "t-online.de" /* T-Mobile */, "web.de", "yahoo.de",

  /* Italian ISP domains */
  "libero.it", "virgilio.it", "hotmail.it", "aol.it", "tiscali.it", "alice.it", "live.it", "yahoo.it", "email.it", "tin.it", "poste.it", "teletu.it",

  /* Russian ISP domains */
  "mail.ru", "rambler.ru", "yandex.ru", "ya.ru", "list.ru",

  /* Belgian ISP domains */
  "hotmail.be", "live.be", "skynet.be", "voo.be", "tvcablenet.be", "telenet.be",

  /* Argentinian ISP domains */
  "hotmail.com.ar", "live.com.ar", "yahoo.com.ar", "fibertel.com.ar", "speedy.com.ar", "arnet.com.ar",

  /* Domains used in Mexico */
  "yahoo.com.mx", "live.com.mx", "hotmail.es", "hotmail.com.mx", "prodigy.net.mx",

  /* Domains used in Brazil */
  "yahoo.com.br", "hotmail.com.br", "outlook.com.br", "uol.com.br", "bol.com.br", "terra.com.br", "ig.com.br", "itelefonica.com.br", "r7.com", "zipmail.com.br", "globo.com", "globomail.com", "oi.com.br"
];