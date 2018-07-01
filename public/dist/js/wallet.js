function Wallet() {

    var componentCounter = 0;

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

    // PRIVATE METHODS

    // Constructor method
    function constructor() {
        if (typeof web3 !== 'undefined') {
            web3js = new Web3(web3.currentProvider);
            if (web3js.eth.coinbase) {
                if (web3js.version.network == networkVersion) {
                    initSuccess();
                } else {
                    initWrongNetwork();
                }                
            } else {
                initNotSignedIn();
            }
        } else {
            initNoMetaMask();
        }
    }

    // Successful init
    function initSuccess() {
        anote = web3js.eth.contract(anoteAbi).at(anoteAddress);

        setValue('address', web3js.eth.coinbase);

        web3js.eth.getBalance(web3js.eth.coinbase, function(error, result) {
            var balance = parseFloat(web3js.fromWei(parseInt(result))).toFixed(5);
            setHTML('balanceEth', balance);
            updateCounter();
        });

        anote.balanceOf(web3js.eth.coinbase, function(error, result){
            var antBalance = parseFloat(web3js.fromWei(parseInt(result))).toFixed(5);
            setHTML('balanceAnt', antBalance);
            updateCounter();
        });
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
                value: web3js.eth.coinbase
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

    // Transfers ETH
    function transferEth(addressTo, amount) {

    }

    // Transfers ANT
    function transferAnt(addressTo, amount) {
        $('#content').fadeOut(function() {
            $('#transactionInProgress').fadeIn();
        });
        amount = web3js.toWei(parseFloat(amount));
        anote.transfer(addressTo, amount, function(err, res) {
            if (err == null) {
                var interval = setInterval(function(){
                    web3js.eth.getTransaction(res, function(err, res) {
                        if (res.blockNumber) {
                            clearInterval(interval);
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
        });
    }

    // Attach all events
    getEl('payButton').addEventListener('click', bind(this, this.pay), false);
    getEl('copyButton').addEventListener('click', bind(this, this.copy), false);

    // Calling Wallet constructor
    constructor();
} 

window.addEventListener('load', function() {
    var wallet = new Wallet();
});
