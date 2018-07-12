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
                case 1:
                    transferEth(addressTo, amount);
                    break;
                default:
                    transferAnt(addressTo, amount);
            }
        }
    }

    // Copy method for copying address to clipboard
    this.copy = function() {
        var copyText = document.getElementById("address");
        copyText.select();
        document.execCommand("copy");
        $('#copymessage').fadeIn(function() {
            setTimeout(() => {
                $('#copymessage').fadeOut();
            }, 2000);
        });
    }

    // Saves user's nickname to blockchain
    this.save = function() {
        var nickname = getEl('nickname').value;
        if (validateNicknameFields(nickname)) {
            anonutopia.setNickname(nickname, function(error, result) {
                console.log(error);
            });
        }
    }

    // Exchange currencies
    this.exchange = function() {
        var selectedCurrency = getEl('currency').selectedIndex;
        var amount = getEl('amount').value;
        var referral = getEl('referral').value;
        if (validateExchangeFields(selectedCurrency, amount, referral)) {
            $('#content').fadeOut(function() {
                $('#transactionInProgress').fadeIn();
            });
            switch (selectedCurrency) {
                case 1:
                    anote.antToEth(web3js.toWei(amount), handleExchangeTransactionResult);
                    break;
                default:
                    anote.ethToAnt(referral, { from: web3js.eth.coinbase, value: web3js.toWei(amount) }, handleExchangeTransactionResult);
            }
        }
    }

    // Withdraw ANT profit
    this.withdraw = function() {
        $('#content').fadeOut(function() {
            $('#transactionInProgress').fadeIn();
        });
        anote.withdrawProfit(handleProfitTransactionResult);
    }

    // Sign in method
    this.signIn = function() {
        var pass = getEl('password').value;
        if (validatePasswordField(pass)) {
            try {
                var restoredPhrase = Waves.Seed.decryptSeedPhrase(Cookies.get('encrypted'), pass);
                Cookies.set('seed', restoredPhrase, { expires: 1 });
                window.location = '/';
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

    // PRIVATE METHODS

    // Constructor method
    function constructor() {
        $(".sidebar-menu a").each(function() {
            if ($(this).attr('href') == window.location.pathname) {
                $(this).parent().addClass('active');
            }
        });

        Waves = WavesAPI.create(WavesAPI.MAINNET_CONFIG);
        var restoredPhrase = Cookies.get('seed');

        if (!restoredPhrase) {
            if (window.location.pathname != '/sign-in/') {
                window.location = '/sign-in/';
                return;
            } else {
                return;
            }
        }

        seed = Waves.Seed.fromExistingPhrase(restoredPhrase);
        Cookies.set('seed', restoredPhrase, { expires: 1 });

        initSuccess();
    }

    // Successful init
    function initSuccess() {
        setValue('address', seed.address);

        Waves.API.Node.v1.assets.balance(seed.address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf").then((balance) => {
            var antBalance = parseFloat(parseFloat(balance.balance) / parseFloat(10**8)).toFixed(5);
            setHTML('balanceAnt', antBalance);
            updateCounter();
        });

        Waves.API.Node.v1.assets.balance(seed.address, "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53").then((balance) => {
            var ethBalance = parseFloat(parseFloat(balance.balance) / parseFloat(10**8)).toFixed(5);
            setHTML('balanceEth', ethBalance);
            updateCounter();
        });

        timeout = setTimeout(initSuccess, 1000);
    }

    // Successful init for profile page
    function initSuccessProfile() {
        anonutopia = web3js.eth.contract(anonutopiaAbi).at(anonutopiaAddress);
        anonutopia.getNickname(function(error, result) {
            if (result.length) {
                setHTML('nicknameTag', result);
                if (!getEl('nickname').value.length) {
                    setValue('nickname', result);
                }
            }
            updateCounter();
        });
        updateCounter();

        timeout = setTimeout(initSuccessProfile, 1000);
    }

    // Successful init for exchange page
    function initSuccessExchange() {
        anote = web3js.eth.contract(anoteAbi).at(anoteAddress);
        anonutopia = web3js.eth.contract(anonutopiaAbi).at(anonutopiaAddress);
        anonutopia.getNickname(function(error, result) {
            if (result.length) {
                setHTML('nicknameTag', result);
            }
        });
        anote.priceBuy(function(error, result) {
            var price = parseFloat(web3js.fromWei(result)).toFixed(5);
            setHTML('priceBuy', price);
            updateCounter();
        });

        timeout = setTimeout(initSuccessExchange, 1000);
    }

    // Successful init for profit page
    function initSuccessProfit() {
        anote = web3js.eth.contract(anoteAbi).at(anoteAddress);
        anonutopia = web3js.eth.contract(anonutopiaAbi).at(anonutopiaAddress);
        anonutopia.getNickname(function(error, result) {
            if (result.length) {
                setHTML('nicknameTag', result);
            }
            updateCounter();
        });
        anote.balanceProfitOf(web3js.eth.coinbase, function(error, result) {
            var balanceProfit = parseFloat(web3js.fromWei(result)).toFixed(5);
            setHTML('profit', balanceProfit);
            updateCounter();
        });

        timeout = setTimeout(initSuccessProfit, 1000);
    }

    // Inits if there's no MetaMask
    function initNoMetaMask() {
        $('#loader').fadeOut(function() {
            $('#nometamask').fadeIn();
        });
    }

    // Inits if user is not signed in
    function initNotSignedIn() {
        $('#loader').fadeOut(function() {
            $('#notsignedin').fadeIn();
        });
    }

    // Inits if user is using the wrong network
    function initWrongNetwork() {
        $('#loader').fadeOut(function() {
            $('#wrongnetwork').fadeIn();
        });
    }

    // Updates counter for loading purposes
    function updateCounter() {
        componentCounter++;
        if (componentCounter == 2) {
            var qr = new QRious({
                size: 300,
                element: document.getElementById('qr'),
                value: 'waves://' + seed.address
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
            setHTML('errorMessagePayment', 'Oba polja su nužna za slanje.');
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

    // Checks and validates fields for nickname form
    function validateNicknameFields(nickname) {
        var validates = true;

        if (nickname.length == 0) {
            $('#nicknameGroup').addClass('has-error');
            validates = false;
        }

        if (!validates) {
            setHTML('errorMessageNickname', 'Polje ne može biti prazno.');
            $('#errorMessageNickname').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessageNickname').fadeOut();
                    $('#nicknameGroup').removeClass('has-error');
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
            setHTML('required', 'Polje ne može biti prazno.');
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
    function validateExchangeFields(selectedCurrency, amount, referral) {
        var validates = true;

        if (amount.length == 0) {
            $('#amountGroup').addClass('has-error');
            validates = false;
        }

        if (!validates) {
            setHTML('errorMessageExchange', 'Polje ne može biti prazno.');
            $('#errorMessageExchange').fadeIn(function() {
                setTimeout(() => {
                    $('#errorMessageExchange').fadeOut();
                    $('#amountGroup').removeClass('has-error');
                }, 2000);
            });
        }

        return validates;
    }

    // Transfers ETH
    function transferEth(addressTo, amount) {
        $('#content').fadeOut(function() {
            $('#transactionInProgress').fadeIn();
        });
        amount = web3js.toWei(parseFloat(amount));
        web3js.eth.sendTransaction({ to: addressTo, value: amount }, handleTransactionResult);
    }

    // Transfers ANT
    function transferAnt(addressTo, amount) {
        // $('#content').fadeOut(function() {
        //     $('#transactionInProgress').fadeIn();
        // });
        // amount = web3js.toWei(parseFloat(amount));
        // anote.transfer(addressTo, amount, handleTransactionResult);
        const transferData = {
            recipient: addressTo,
            assetId: '4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf',
            amount: amount * 10**8,
            feeAssetId: 'WAVES',
            fee: 100000,
            attachment: '',
            timestamp: Date.now()
        };

        Waves.API.Node.v1.assets.transfer(transferData, seed.keyPair).then((responseData) => {
            console.log(responseData);
        });
    }

    // Handles transaction result
    function handleTransactionResult(err, res) {
        if (err == null) {
            var interval = setInterval(function(){
                web3js.eth.getTransaction(res, function(err, res) {
                    if (res.blockNumber) {
                        clearInterval(interval);
                        clearTimeout(timeout);
                        initSuccess();
                        $('#transactionInProgress').fadeOut(function() {
                            $('#transactionSuccess').fadeIn(function() {
                                setTimeout(() => {
                                    $('#transactionSuccess').fadeOut();
                                }, 10000);
                            });
                            $('#content').fadeIn();
                        });
                    }
                });
            }, 1000);
        } else {
            $('#transactionInProgress').fadeOut(function() {
                $('#transactionError').fadeIn(function() {
                    setTimeout(() => {
                        $('#transactionError').fadeOut();
                    }, 10000);
                });
                $('#content').fadeIn();
            });
        }
    }

    // Handles exchange transaction result
    function handleExchangeTransactionResult(error, result) {
        if (error == null) {
            var interval = setInterval(function(){
                web3js.eth.getTransaction(result, function(err, res) {
                    if (res.blockNumber) {
                        clearInterval(interval);
                        clearTimeout(timeout);
                        initSuccessExchange();
                        $('#transactionInProgress').fadeOut(function() {
                            $('#transactionSuccess').fadeIn(function() {
                                setTimeout(() => {
                                    $('#transactionSuccess').fadeOut();
                                }, 10000);
                            });
                            $('#content').fadeIn();
                        });
                    }
                });
            }, 1000);
        } else {
            $('#transactionInProgress').fadeOut(function() {
                $('#transactionError').fadeIn(function() {
                    setTimeout(() => {
                        $('#transactionError').fadeOut();
                    }, 10000);
                });
                $('#content').fadeIn();
            });
        }
    }

    // Handles profit transaction result
    function handleProfitTransactionResult(error, result) {
        if (error == null) {
            var interval = setInterval(function(){
                web3js.eth.getTransaction(result, function(err, res) {
                    if (res.blockNumber) {
                        clearInterval(interval);
                        clearTimeout(timeout);
                        initSuccessProfit();
                        $('#transactionInProgress').fadeOut(function() {
                            $('#transactionSuccess').fadeIn(function() {
                                setTimeout(() => {
                                    $('#transactionSuccess').fadeOut();
                                }, 10000);
                            });
                            $('#content').fadeIn();
                        });
                    }
                });
            }, 1000);
        } else {
            $('#transactionInProgress').fadeOut(function() {
                $('#transactionError').fadeIn(function() {
                    setTimeout(() => {
                        $('#transactionError').fadeOut();
                    }, 10000);
                });
                $('#content').fadeIn();
            });
        }
    }

    // Handling currency select changed state
    function exchangeCurrencyChanged() {
        var selectedCurrency = getEl('currency').selectedIndex;
        if (selectedCurrency == 0) {
            $('#referralGroup').fadeIn();
        } else {
            $('#referralGroup').fadeOut();
        }
    }

    // Attach all events
    switch (window.location.pathname) {
        case '/':
            getEl('payButton').addEventListener('click', bind(this, this.pay), false);
            getEl('copyButton').addEventListener('click', bind(this, this.copy), false);
            break;
        case '/profit/':
            getEl('withdrawButton').addEventListener('click', bind(this, this.withdraw), false);
            break;
        case '/profile/':
            getEl('saveButton').addEventListener('click', bind(this, this.save), false);
            break;
        case '/exchange/':
            getEl('exchangeButton').addEventListener('click', bind(this, this.exchange), false);
            getEl('currency').addEventListener('change', bind(this, exchangeCurrencyChanged), false);
            break;
        case '/sign-in/':
            getEl('signInButton').addEventListener('click', bind(this, this.signIn), false);
            break;
    }

    // Calling Wallet constructor
    constructor();
}
