var BASE_URL = "http://" + document.location.hostname + ":8080";

var initValue = function(name) {
  $.get(BASE_URL + "/module/" + name + "/properties", function(data) {
    for (var k in data) {
      $("#" + name + "-" + k).val(data[k]);
    }
  });
}

var loadModule = function(name) {
  $.get(BASE_URL + "/module/" + name + "/template", function(data) {
    $("#setting-templates").append(data);
    initValue(name);
  });
}

var loadTemplates = function() {
  $.get(BASE_URL + "/module", function(data) {
    for (var i = 0; i < data.length; i++) {
      loadModule(data[i]);
    }
  })
}

var saveProperties = function(module) {
  var params = $('#' + module + '-form').serializeArray();

  var data = {};
  data.name = module;
  data.properties = {};

  for (var i in params) {
    var name = params[i].name.replace(module + '-', "");
    if ($.isNumeric(params[i].value)) {
      data.properties[name] = parseInt(params[i].value);
    } else {
      data.properties[name] = params[i].value;
    }
  }

  console.log(JSON.stringify(data));

  $.post(BASE_URL + "/module/" + module + "/properties", JSON.stringify(data), function(data) {
    console.log(data);
  });
  return false;
}

loadTemplates();