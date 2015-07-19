
/* batsocket conn */

(function() {

	//'use strict';

	if (!window["WebSocket"]) {
		console.log('Batbrowser not identified!! Access denied!!');
		return false;
	}

	var sock = new SockJS('http://localhost:8000/batcave');

	sock.onopen = function() {
		console.log('open');
	};

	sock.onmessage = function(e) {
		console.log('message', e.data);
	};

	sock.onclose = function() {
		console.log('close');
	};

	//sock.send('test');
	sock.close();
}());
