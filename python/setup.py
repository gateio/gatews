from setuptools import find_packages, setup

VERSION = '0.3.1'

setup(
    name='gate-ws',
    version=VERSION,
    packages=find_packages(),
    url='https://github.com/gateio/gatews',
    install_requires=['websockets>=8.1'],
    license='MIT License',
    author='gateio',
    keywords=["Gate WebSocket V4"],
    author_email='dev@mail.gate.io',
    description='Gate.io WebSocket V4 Python SDK'
)
