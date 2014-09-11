(function () {

	function fromNow() {
		return function(date) {
			return moment(new Date(date*1000)).fromNow();
		}
	}

	function toDuration() {
		return function(seconds) {
			return moment.duration(seconds, "seconds").humanize();
		}
	}

	function toDate() {
		return function(date) {
			return moment(new Date(date*1000)).format('ll');
		}
	}

	function pubLink() {
		return function(pkg) {
			if (pkg === undefined || pkg.name === undefined) {
				return;
			}
			return 'https://pub.dartlang.org/packages/'+pkg.name;
		}
	}

	angular
		.module('app')
		.filter('pubLink', pubLink)
		.filter('toDate', toDate)
		.filter('toDuration', toDuration)
		.filter('fromNow', fromNow);

})();