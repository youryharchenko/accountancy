// Transaction (trait) - script

function constructor(trait, obj) {
    //var traitStr = JSON.stringify(trait)
    //var objStr = JSON.stringify(obj)
    var result = {};
    if (obj.pay) {
        switch (obj.pay.request.service) {
            case "megapolissto":
            case "nparts":
                result.transactionId = obj.pay.request.transactionId;
                result.terminalId = obj.pay.request.terminalId;
                result.service = obj.pay.request.service;
                result.status = obj.pay.response.status;
                result.invoiceId = obj.pay.request.body.invoiceId;
                result.method = obj.pay.request.body.method;
                if (result.method == 'cash') result.pushAmount = obj.pay.request.body.pushAmount;
                else result.pushAmount = 0;
                result.amount = obj.pay.request.body.amount;
                result.payDate = obj.pay.request.body.payDate;
                result.purpose = obj.pay.request.body.purpose;
                result.entrStat = obj.entrStat;
                result.updAt = obj.updAt;
                break;
        }
    }

    return result;
}