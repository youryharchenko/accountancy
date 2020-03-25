// Transaction (trait) - script

function constructor(trait, obj) {
    //var traitStr = JSON.stringify(trait)
    //var objStr = JSON.stringify(obj)
    var result = {};
    if (obj.pay) {
        result.transactionId = obj.pay.request.transactionId;
        result.terminalId = obj.pay.request.terminalId;
        result.service = obj.pay.request.service;
        result.status = obj.pay.response.status;
        result.body = obj.pay.request.body;
        result.updAt = obj.updAt;
    }

    return result;
}