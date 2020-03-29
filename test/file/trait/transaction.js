// Transaction (trait) - script

function constructor(trait, obj) {
    //var traitStr = JSON.stringify(trait)
    var props = {};
    var status = -1.0;
    var hash = "";
    var batch = [];
    if (obj.pay) {
        switch (obj.pay.request.service) {
            case "megapolissto":
            case "nparts":

                //log(obj.pay.request.transactionId);

                props.transactionId = obj.pay.request.transactionId;
                props.terminalId = obj.pay.request.terminalId;
                props.service = obj.pay.request.service;
                props.status = obj.pay.response.status;
                props.invoiceId = obj.pay.request.body.invoiceId;
                props.method = obj.pay.request.body.method;
                if (props.method == 'cash') props.pushAmount = obj.pay.request.body.pushAmount;
                else props.pushAmount = 0;
                props.amount = obj.pay.request.body.amount;
                props.payDate = obj.pay.request.body.payDate;
                props.purpose = obj.pay.request.body.purpose;
                props.entrStat = obj.entrStat;
                props.updAt = obj.updAt;

                hash = makeHash(props);

                var model = obj.pay.request.model;
                if (model) {
                    var hashModel = makeHash(model);

                    var insModel = {
                        "request": {
                            "command": "insert",
                            "service": "object"
                        },
                        "body": {
                            "name": hashModel,
                            "props": model,
                            "traits": ["Model"]
                        }
                    };

                    batch.push(insModel);
                }

                status = 0.0;
                break;
        }
    }

    return {
        "props": props,
        "hash": hash,
        "batch": batch,
        "status": status
    };
}