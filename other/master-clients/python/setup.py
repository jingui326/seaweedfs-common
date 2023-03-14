# coding: utf-8

"""
    Seaweedfs Master Server API

    The Seaweedfs Master Server API allows you to store blobs  # noqa: E501

    The version of the OpenAPI document: 3.43.0
    Generated by: https://openapi-generator.tech
"""


from setuptools import setup, find_packages  # noqa: H301

# To install the library, run the following
#
# python setup.py install
#
# prerequisite: setuptools
# http://pypi.python.org/pypi/setuptools
NAME = "openapi-client"
VERSION = "1.0.0"
PYTHON_REQUIRES = ">=3.7"
REQUIRES = [
    "urllib3 >= 1.25.3",
    "python-dateutil",
    "pydantic",
    "aenum"
]

setup(
    name=NAME,
    version=VERSION,
    description="Seaweedfs Master Server API",
    author="OpenAPI Generator community",
    author_email="team@openapitools.org",
    url="",
    keywords=["OpenAPI", "OpenAPI-Generator", "Seaweedfs Master Server API"],
    install_requires=REQUIRES,
    packages=find_packages(exclude=["test", "tests"]),
    include_package_data=True,
    long_description_content_type='text/markdown',
    long_description="""\
    The Seaweedfs Master Server API allows you to store blobs  # noqa: E501
    """
)
