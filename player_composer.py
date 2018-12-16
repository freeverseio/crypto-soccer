#!/usr/bin/env python
import os, sys
import subprocess
from xml.dom import minidom
import argparse
import sha3
from playerdb import *

# see : /System/Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/xml/dom/minidom.py

def get_decimal_hash(x):
    return int(sha3.sha3_256(x).hexdigest(),16)

def rgb2hex(r,g,b):
    hex = "#{:02x}{:02x}{:02x}".format(r,g,b)
    return hex

def get_svg_header(x=0, y=0, w=700, h=1000):
    return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="{x} {y} {w} {h}" xmlns:xlink="http://www.w3.org/1999/xlink">'.format(x=x,y=y,w=w,h=h)

def get_svg_footer():
    return '</svg>'

class SvgPath:
    fill=None
    d=None
    stroke=None
    stroke_width=None
    stroke_linecap=None
    stroke_linejoin=None

    def __init__(self, fill, d, stroke=None, stroke_width=None, stroke_linecap=None, stroke_linejoin=None):
        self.fill = fill
        self.d = d
        self.stroke=stroke
        self.stroke_width=stroke_width
        self.stroke_linecap=stroke_linecap
        self.stroke_linejoin=stroke_linejoin

    def toxml(self):
        result='\n<path fill="{p.fill}" d="{p.d}"'.format(p=self)

        if self.stroke:
            result += ' stroke="{p.stroke}"'.format(p=self)
        if self.stroke_width:
            result += ' stroke-width="{p.stroke_width}"'.format(p=self)
        if self.stroke_linecap:
            result += ' stroke-linecap="{p.stroke_linecap}"'.format(p=self)
        if self.stroke_linejoin:
            result += ' stroke-join="{p.stroke_linejoin}"'.format(p=self)

        result += '/>'

        return result

class SvgGroup:
    name=None
    transform=None
    paths=[]
    groups=[]
    display=None

    def __init__(self):
        pass

    def __init__(self, name=None, transform=None, paths=[], groups=[], display=None):
        self.name = name
        self.transform = transform
        self.display = display

        if paths:
            if not isinstance(paths, SvgPath) and not isinstance(paths, list):
                print 'Error: paths must be either an SvgPath or a list of SvgPath'
                sys.exit(-1)
            if isinstance(paths,SvgPath):
                self.paths = [paths]
            else:
                self.paths = paths
        else:
            self.paths = []

        if groups:
            if not isinstance(groups, SvgGroup) and not isinstance(groups, list):
                print 'Error: groups must be either an SvgGroup or a list of SvgGroup'
                sys.exit(-1)
            if isinstance(groups,SvgGroup):
                self.groups = [groups]
            else:
                self.groups = groups
        else:
            self.groups = []


    def toxml(self):
        result='\n<g'

        if self.name:
            result += ' id="{p.name}"'.format(p=self)
        if self.transform:
            result += ' transform="{p.transform}"'.format(p=self)
        if self.display:
            result += ' display="{p.display}"'.format(p=self)
        result+='>'

        for path in self.paths:
            result += path.toxml()
        for group in self.groups:
            result += group.toxml()

        result+='\n</g>'
        return result

def create_svg(groups, filename):
    f = open(filename + '.svg', 'w')
    f.write(get_svg_header())
    f.write('\n<g>')

    for group in groups:
        f.write(group.toxml())

    f.write('\n</g>')

    f.write(get_svg_footer())
    f.close()

def get_arms(n, color):
    return SvgPath(
        fill=color,
        d=arms_db[n],
        )

def get_head(n, color):
    return SvgPath(
        fill=color,
        d=head_db[n],
        )
def get_neck(n, color):
    return SvgPath(
        fill=color,
        d=neck_db[n],
        )

def get_body(n, color):
    return SvgGroup(
        name = 'body',
        transform = "matrix( 1, 0, 0, 1, 0,0)",
        paths = [get_head(n,color), get_neck(n, color), get_arms(n, color)]
        )

def get_body_color(hash_str):
    return '#' + hash_str[0:6]

def generate_player(name):
    hash_str = sha3.sha3_256(name).hexdigest()
    body_style = 0
    body_color = get_body_color(hash_str)

    return [
            get_body(body_style, body_color),
           ]


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='script to extract elements from svg or merge svgs',
            formatter_class=argparse.ArgumentDefaultsHelpFormatter,
            usage='./generate_player -n <name> -o <output_name>'
            )
    parser.add_argument('-n', '--name', help='name from which to generate a player', required=True)
    parser.add_argument('-o', '--output', help='name from which to generate a player', required=True)

    args = parser.parse_args()

    player_name = args.name
    output_name = args.output

    create_svg(generate_player(player_name), output_name)
