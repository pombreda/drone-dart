'use script';

(function () {

	/**
	 * Creates the angular application.
	 */
	angular.module('app', [
			'ngRoute',
			'ui.filters'
		]);

	/**
	 * Defines the route configuration for the
	 * main application.
	 */
	function Config ($routeProvider) {
		$routeProvider
		.when('/:package/:version/channel/:channel/sdk/:sdk', {
			templateUrl: '/static/scripts/views/build.html',
			controller: 'BuildCtrl'
		})
		.when('/:package/:version/channel/:channel', {
			templateUrl: '/static/scripts/views/build.html',
			controller: 'BuildCtrl'
		})
		.when('/:package/:version', {
			templateUrl: '/static/scripts/views/build.html',
			controller: 'BuildCtrl'
		})
		.when('/', {
			templateUrl: '/static/scripts/views/main.html',
			controller: 'MainCtrl'
		});
	}

	angular
		.module('app')
		.config(Config);

})();