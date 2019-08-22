#! /usr/bin/env bash

set -euo pipefail

OS="$(go env GOHOSTOS)"
ARCH="$(go env GOARCH)"

# creates links to generated go protos

echo -e ">>> Compiling the Go proto"
for label in $(bazel query 'kind(go_proto_library, //...)'); do
	package="${label%%:*}"
	package="${package##//}"
	target="${label##*:}"

	# do not continue if the package does not exist
	[[ -d "${package}" ]] || continue

	# compute the path where bazel put the files
	out_path="bazel-bin/${package}/${OS}_${ARCH}_stripped/${target}%/miniboard.app/${package}"

    echo $out_path

	# compute the relative_path to the
	count_paths="$(echo -n "${package}" | tr '/' '\n' | wc -l)"
	relative_path=""
	for i in $(seq 0 ${count_paths}); do
		relative_path="../${relative_path}"
	done

	bazel build "${label}"

	found=0
	for f in ${out_path}/*.pb.go; do
		if [[ -f "${f}" ]]; then
			found=1
			ln -nsf "${relative_path}${f}" "${package}/"
		fi
	done
	if [[ "${found}" == "0" ]]; then
		echo "ERR: no .pb.go file was found inside $out_path for the package ${package}"
		exit 1
	fi

    found=0
	for f in ${out_path}/*.gw.go; do
		if [[ -f "${f}" ]]; then
			found=1
			ln -nsf "${relative_path}${f}" "${package}/"
		fi
	done
	if [[ "${found}" == "0" ]]; then
		echo "ERR: no .gw.go file was found inside $out_path for the package ${package}"
		exit 1
	fi

done
