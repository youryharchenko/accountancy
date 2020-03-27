// Model (trait) - script

function constructor(trait, obj) {
    //var traitStr = JSON.stringify(trait)
    //var objStr = JSON.stringify(obj)
    var result = {};
    if (obj.pay && obj.pay.request.model) {
        result = obj.pay.request.model;
    }

    return result;
}