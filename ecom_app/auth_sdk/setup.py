from setuptools import setup, find_packages

with open('requirements.txt') as f:
    requirements = f.read().splitlines()

setup(
    name='authsdk',
    version='0.1',
    packages=['authsdk'],
    install_requires=requirements,
)
