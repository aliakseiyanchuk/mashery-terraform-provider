resource "mashery_processor_chain" "example" {
  processor_chain {
    adapter = "com.acme.setHeader.36"
    pre_process_enabled = true
    pre_config = [
      "b",
      "c",
      "d"]
    post_process_enabled = true
    post_config = [
      "d",
      "e",
      "f"]
  }
  processor_chain {
    adapter = "doSomethingElseAdapter2.cc.dd1.5"
    pre_process_enabled = true
    pre_config = [
      "b1",
      "b12",
      "c1",
      "d1"]
    post_process_enabled = true
    post_config = [
      "d1",
      "1e",
      "f1"]
  }
  processor_chain {
    adapter = "andEventMorre"
    pre_process_enabled = true
    pre_config = [
      "g1",
      "h12",
      "i1",
      "j1"]
    post_process_enabled = true
    post_config = [
      "k1",
      "le",
      "m1"]
  }
}

output "compiled_inputs" {
  value = mashery_processor_chain.example.pre_config
}

output "compiled_chain" {
  value = mashery_processor_chain.example.compiled_chain
  depends_on = [mashery_processor_chain.example]
}