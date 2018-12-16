#!/usr/bin/env python
import os, sys
import subprocess
from xml.dom import minidom
import argparse

# see : /System/Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/xml/dom/minidom.py

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

def createSvg(groups, filename):
    f = open(filename + '.svg', 'w')
    f.write(get_svg_header())
    f.write('\n<g>')

    for group in groups:
        f.write(group.toxml())

    f.write('\n</g>')

    f.write(get_svg_footer())
    f.close()

def generatePlayer():
    result =''
    return result;

arms_path = SvgPath(
        fill = '#33CC66',
        d = 'M471.9,626.25l-5.45-20.3c-0.446-1.16-0.862-2.36-1.25-3.601l22.7,84.75L404,750.95 c-10.797,8.253-17.13,19.103-19,32.55c-1.867,13.452,1.316,25.619,9.55,36.5c8.26,10.87,19.109,17.236,32.55,19.1 c13.445,1.873,25.646-1.311,36.601-9.55L572.85,746.5c0.876-0.613,1.709-1.264,2.5-1.95l4.2-3.8 c5.813-5.824,9.83-12.708,12.05-20.65c1.723-5.592,2.34-11.325,1.851-17.199c-0.231-4.337-1.064-8.537-2.5-12.601 l-22.601-84.35c1.807,11.115-0.21,21.749-6.05,31.899c-6.802,11.829-16.751,19.529-29.85,23.101 c-13.104,3.572-25.571,1.955-37.4-4.851C483.223,649.297,475.506,639.347,471.9,626.25 M137.55,637.8 c-5.839-10.15-7.873-20.784-6.1-31.899l-22.55,84.25c-1.454,4.131-2.304,8.398-2.55,12.8c-0.438,5.849,0.195,11.549,1.9,17.1 c2.147,7.642,5.963,14.309,11.45,20l4.3,4.101c0.991,0.826,2.007,1.609,3.05,2.35l109.05,82.95 c10.865,8.238,23.015,11.438,36.45,9.6c13.468-1.804,24.334-8.138,32.6-19c8.303-10.83,11.537-22.98,9.7-36.45 c-1.805-13.436-8.139-24.336-19-32.699L211.9,687.05l22.7-84.75c-0.346,1.212-0.746,2.395-1.2,3.55l-5.5,20.351 c-3.572,13.097-11.272,23.046-23.1,29.85c-11.829,6.806-24.296,8.422-37.4,4.851C154.302,657.329,144.352,649.629,137.55,637.8z'
        )

arms_group = SvgGroup(
    name = "Brazos",
    transform = "matrix( 1, 0, 0, 1, 0,0)",
    paths = [arms_path])



if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='script to extract elements from svg or merge svgs',
            formatter_class=argparse.ArgumentDefaultsHelpFormatter,
            usage=''
            )
    parser.add_argument('-n', '--name', help='name from which to generate a player', required=True)
    parser.add_argument('-o', '--output', help='name from which to generate a player', required=True)

    args = parser.parse_args()

    player_name = args.name
    output_name = args.output

    createSvg([arms_group], output_name)
