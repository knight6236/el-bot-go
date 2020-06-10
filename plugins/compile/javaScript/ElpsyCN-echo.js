var whenDetail = {
    type: "Plain",
    regex: "echo\\s(.+)"
};
var doDetail = {
    type: "Plain",
    text: "{el-regex-0}"
};

var whenMessage = {
    detail: [whenDetail]
};
var doMessage = {
    quote: true,
    detail: [doDetail]
};

var config = {
    global: [{ when: { message: whenMessage }, do: { message: doMessage } }]
};



var args = process.argv.splice(2);
var jsonStr = args[0];
jonsObj = JSON.parse(jsonStr);
if (jonsObj.enable == true) {
    var jsonStr = JSON.stringify(config);
    process.stdout.write(jsonStr);
} else {
    process.stdout.write("");
}