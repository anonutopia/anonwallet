function Wallet() {

    var componentCounter = 0;

    // Payment method
    this.pay = function() {
        console.log('payment');
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
        var anote = web3js.eth.contract(anoteAbi).at(anoteAddress);

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

    // Attach all events
    getEl('payButton').addEventListener('click', bind(this, this.pay), false);

    // Calling Wallet constructor
    constructor();
} 

window.addEventListener('load', function() {
    var wallet = new Wallet();
});
