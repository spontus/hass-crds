#!/usr/bin/env python3
"""Generate Kubernetes CRDs for Home Assistant MQTT entities."""

import os
import sys
from pathlib import Path

import yaml

# Add the generator directory to the path
sys.path.insert(0, str(Path(__file__).parent))

from schemas.common import STATUS_SCHEMA, get_all_common_properties
from schemas.entities import ALL_ENTITIES

API_GROUP = "mqtt.home-assistant.io"
API_VERSION = "v1alpha1"


def build_spec_schema(entity: dict, include_common: bool = True) -> dict:
    """Build the OpenAPI schema for a CRD's spec."""
    properties = {}

    # Add entity-specific properties
    for name, prop in entity["properties"].items():
        properties[name] = prop.copy()

    # Add common properties for entity types (not for MQTTDevice)
    if include_common and entity["component"] is not None:
        for name, prop in get_all_common_properties().items():
            properties[name] = prop.copy()

    schema = {
        "type": "object",
        "properties": properties,
    }

    # Add required fields
    if entity["required"]:
        schema["required"] = entity["required"]

    return schema


def build_printer_columns(entity: dict) -> list:
    """Build printer columns for kubectl get output."""
    columns = [
        {
            "name": "Name",
            "type": "string",
            "description": "Display name in Home Assistant",
            "jsonPath": ".spec.name",
        },
        {
            "name": "Published",
            "type": "string",
            "description": "Whether discovery has been published",
            "jsonPath": ".status.conditions[?(@.type=='Published')].status",
        },
        {
            "name": "Last Published",
            "type": "date",
            "description": "When discovery was last published",
            "jsonPath": ".status.lastPublished",
        },
        {
            "name": "Age",
            "type": "date",
            "jsonPath": ".metadata.creationTimestamp",
        },
    ]
    return columns


def build_crd(entity: dict) -> dict:
    """Build a complete CRD definition for an entity type."""
    # Determine if this is a special utility type (MQTTDevice)
    is_utility = entity["component"] is None

    crd = {
        "apiVersion": "apiextensions.k8s.io/v1",
        "kind": "CustomResourceDefinition",
        "metadata": {
            "name": f"{entity['plural']}.{API_GROUP}",
            "annotations": {
                "controller-gen.kubebuilder.io/version": "v0.14.0",
            },
            "labels": {
                "app.kubernetes.io/name": "hass-crds",
                "app.kubernetes.io/component": "crds",
            },
        },
        "spec": {
            "group": API_GROUP,
            "names": {
                "kind": entity["kind"],
                "listKind": f"{entity['kind']}List",
                "plural": entity["plural"],
                "singular": entity["singular"],
            },
            "scope": "Namespaced",
            "versions": [
                {
                    "name": API_VERSION,
                    "served": True,
                    "storage": True,
                    "additionalPrinterColumns": build_printer_columns(entity),
                    "schema": {
                        "openAPIV3Schema": {
                            "type": "object",
                            "description": entity["description"],
                            "properties": {
                                "apiVersion": {
                                    "type": "string",
                                    "description": "APIVersion defines the versioned schema of this representation of an object",
                                },
                                "kind": {
                                    "type": "string",
                                    "description": "Kind is a string value representing the REST resource this object represents",
                                },
                                "metadata": {
                                    "type": "object",
                                },
                                "spec": build_spec_schema(entity, include_common=not is_utility),
                                "status": STATUS_SCHEMA,
                            },
                        },
                    },
                    "subresources": {
                        "status": {},
                    },
                },
            ],
        },
    }

    # Add short names if defined
    if entity.get("short_names"):
        crd["spec"]["names"]["shortNames"] = entity["short_names"]

    # Add categories
    crd["spec"]["names"]["categories"] = ["hass", "mqtt"]

    return crd


def generate_crds(output_dir: Path) -> list:
    """Generate all CRD YAML files."""
    generated_files = []

    for entity in ALL_ENTITIES:
        crd = build_crd(entity)
        filename = f"{entity['singular']}.yaml"
        filepath = output_dir / filename

        with open(filepath, "w") as f:
            yaml.dump(crd, f, default_flow_style=False, sort_keys=False, allow_unicode=True)

        generated_files.append(filepath)
        print(f"Generated: {filepath}")

    return generated_files


def concatenate_crds(files: list, output_file: Path) -> None:
    """Concatenate all CRD files into a single file."""
    with open(output_file, "w") as outf:
        outf.write("# Generated CRDs for hass-crds\n")
        outf.write(f"# API Group: {API_GROUP}\n")
        outf.write(f"# Version: {API_VERSION}\n")
        outf.write(f"# Total CRDs: {len(files)}\n")
        outf.write("#\n")
        outf.write("# Install with: kubectl apply -f crds.yaml\n")
        outf.write("# Verify with: kubectl get crds | grep mqtt.home-assistant.io\n")
        outf.write("---\n")

        for i, filepath in enumerate(sorted(files)):
            with open(filepath) as inf:
                content = inf.read()
                outf.write(content)
                if i < len(files) - 1:
                    outf.write("---\n")

    print(f"\nCombined CRDs written to: {output_file}")


def main():
    # Determine paths
    script_dir = Path(__file__).parent
    crd_dir = script_dir.parent
    bases_dir = crd_dir / "bases"
    combined_file = crd_dir / "crds.yaml"

    # Ensure output directory exists
    bases_dir.mkdir(parents=True, exist_ok=True)

    print(f"Generating {len(ALL_ENTITIES)} CRDs...")
    print(f"Output directory: {bases_dir}")
    print()

    # Generate individual CRD files
    generated_files = generate_crds(bases_dir)

    # Create combined file
    concatenate_crds(generated_files, combined_file)

    print(f"\nSuccessfully generated {len(generated_files)} CRDs!")
    print(f"\nTo install: kubectl apply -f {combined_file}")
    print(f"To verify:  kubectl get crds | grep {API_GROUP} | wc -l  # Should be {len(generated_files)}")


if __name__ == "__main__":
    main()
