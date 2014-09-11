(function () {

	/**
	 * Build controller responsible for displaying
	 * the build page and rendering the build results.
	 */	
	function BuildCtrl($scope, $routeParams, packages, builds) {
		var package = $routeParams.package;
		var version = $routeParams.version;

		// gets the build data from the server
		builds.get(package, version).then(function(build){
			$scope.build = build.data;
		}).catch(function(error){
			$scope.error = error;
		});

		// gets the package data from the server
		packages.get(package, version).then(function(pkg){
			$scope.pkg = pkg.data;
		}).catch(function(error){
			$scope.error = error;
		});

		// gets the build output from the server
		builds.getOutput(package, version).then(function(output){
			$scope.output = output.data;
		}).catch(function(error){
			$scope.error = error;
		});
	}

	angular
		.module('app')
		.controller('BuildCtrl', BuildCtrl);
})();