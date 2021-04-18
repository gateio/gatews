from setuptools import find_packages, setup

REQUIRES = [line.strip() for line in open('requirements.txt') if not line.startswith('#')]
VERSION = '0.1.1'

setup(
    name='gate-ws',
    version=VERSION,
    packages=find_packages(),
    url='https://github.com/gateio/gatews',
    install_requires=REQUIRES,
    license='MIT License',
    author='gateio',
    keywords=["Gate WebSocket V4"],
    author_email='dev@mail.gate.io',
    description='Gate.io WebSocket V4 Python SDK'
)
