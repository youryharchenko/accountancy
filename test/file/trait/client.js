// Model (trait) - script

function constructor(trait, obj) {
    //var traitStr = JSON.stringify(trait)
    //var objStr = JSON.stringify(obj)
    var props = {};
    var status = -1.0;
    var hash = "";
    var batch = [];

    props = obj;
    hash = makeHash(props);
    status = 0.0;
    

    return {
        "props": props,
        "hash": hash,
        "batch": {
            "batch": batch
        },
        "status": status
    };
}