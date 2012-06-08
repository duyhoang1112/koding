<?php

require_once 'helpers.php';

class Kite {
  public $kite_name;
  public $uri;
  
  public function __construct ($kite_name, $uri) {
    $this->kite_name = $kite_name;
    $this->uri = $uri;
  }
  
  public function __call ($method, $arguments) {
    if (!isset($this->$method)) {
      $session = get_session();
      $username = isset($session['username']) ? $session['username'] : 'kc';
      $args = array(
        'data' => json_encode(array(
          'toDo'      => $method,
          'withArgs'  => $arguments[0],
          'username'  => $username,          
        )),
      );
      error_log('connecting to '.$this->uri);
      $res = @file_get_contents($this->uri.'?'.http_build_query($args));
      error_log('got a response: '.$res);
      return isset($res) ? $res : array('error' => 503, 'uri' => $uri);
    }
  }
  
  public function __toString () {
    return $this->uri;
  }
}