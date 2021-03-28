package goebaykleinanzeigen

// ParamName represents a parameter name for a category specific parameter
type ParamName string

const (
	CarManufacturer       ParamName = "autos.marke_s"
	CarModel              ParamName = "autos.model_s"
	CarKM                 ParamName = "autos.km_i"
	CarYearOfRegistration ParamName = "autos.ez_i"
	CarHP                 ParamName = "autos.power_i"
	CarTUEV               ParamName = "autos.tuevy_i"
)
