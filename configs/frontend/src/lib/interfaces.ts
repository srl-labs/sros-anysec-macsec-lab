export interface FormInput {
  id:      string;
  label:   string;
  min:     number;
  max:     number;
  step:    number;
  default: number;
}

export interface ServiceState {
  vll:  boolean
  vpls: boolean
  vprn: boolean
}

export interface LinkState {
  bottom:  boolean
  top:     boolean
}

export interface AllState {
  icmp:   ServiceState
  link:   LinkState
  anysec: ServiceState
}

export interface PageData {
  urlHost:   string
  baseUrl:   string
  fetchUrl:  string
  state:     AllState
}