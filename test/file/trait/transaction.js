// Transaction (trait) - script
function makeInsObject(command, service, name, props, traits) {
    return {
        "request": {
            "command": command,
            "service": service
        },
        "body": {
            "name": name,
            "props": props,
            "traits": traits
        }
    }
}

function makeInsLink(command, service, relation, traitFrom, traitTo, objectFrom, objectTo, props) {
    return {
        "request": {
            "command": command,
            "service": service
        },
        "body": {
            "relation": relation,
            "trait-from": traitFrom,
            "trait-to": traitTo,
            "object-from": objectFrom,
            "object-to": objectTo,
            "props": props
        }
    }
}

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
                    var insObject = makeInsObject("insert", "object", hashModel, model, ["Model"]);
                    batch.push(insObject);

                    var insLink = makeInsLink("insert", "link", "ref-to-model", "Transaction", "Model", props.transactionId, hashModel, {})
                    batch.push(insLink);
                }




                var clientName = obj.pay.request.body.clientName;
                if (clientName && clientName.length > 0) {
                    var client = {
                        "clientName": clientName,
                        "clientAddress": obj.pay.request.body.clientAddress, 
                        "clientBirthday": obj.pay.request.body.clientBirthday, 
                        "clientDocument": obj.pay.request.body.clientDocument, 
                        "clientIssueDate": obj.pay.request.body.clientIssueDate, 
                        "clientIssuedBy": obj.pay.request.body.clientIssuedBy, 
                        "clientTaxNumber": obj.pay.request.body.clientTaxNumber
                    };
                    //var hashClient = makeHash(client);
                    var insObject = makeInsObject("insert", "object", clientName, client, ["Client"]);
                    batch.push(insObject);

                    var insLink = makeInsLink("insert", "link", "ref-to-client", "Transaction", "Client", props.transactionId, hashModel, {})
                    batch.push(insLink);
                }

                status = 0.0;
                break;
        }
    }

    return {
        "props": props,
        "hash": hash,
        "batch": {
            "batch": batch
        },
        "status": status
    };
}