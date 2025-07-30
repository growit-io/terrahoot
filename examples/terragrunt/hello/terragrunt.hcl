locals {
  root_dir = get_repo_root()

  terraform_source = "${local.root_dir}/examples/terraform/hello"

  __terraform_source_files = [
    for file in split("\n", run_cmd(
      "--terragrunt-quiet",
      "terrahoot",
      "module-files",
      local.terraform_source == "" ? "." : local.terraform_source
    )) : mark_as_read(file)
  ]

  __tool_version_files = [
    for file in [
      "${local.root_dir}/.terraform-version",
      "${local.root_dir}/.terragrunt-version"
    ] : mark_as_read(file)
  ]
}

terraform {
  source = local.terraform_source
}
