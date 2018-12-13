#!/usr/bin/env python
import os, sys
import subprocess
from xml.dom import minidom
import argparse

# see : /System/Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/xml/dom/minidom.py

svg_header = '''<?xml version="1.0" encoding="utf-8"?>
<!-- Generator: Adobe Illustrator 16.0.0, SVG Export Plug-In . SVG Version: 6.00 Build 0)  -->
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd" [
	<!ENTITY ns_extend "http://ns.adobe.com/Extensibility/1.0/">
	<!ENTITY ns_ai "http://ns.adobe.com/AdobeIllustrator/10.0/">
	<!ENTITY ns_graphs "http://ns.adobe.com/Graphs/1.0/">
	<!ENTITY ns_vars "http://ns.adobe.com/Variables/1.0/">
	<!ENTITY ns_imrep "http://ns.adobe.com/ImageReplacement/1.0/">
	<!ENTITY ns_sfw "http://ns.adobe.com/SaveForWeb/1.0/">
	<!ENTITY ns_custom "http://ns.adobe.com/GenericCustomNamespace/1.0/">
	<!ENTITY ns_adobe_xpath "http://ns.adobe.com/XPath/1.0/">
]>
<svg version="1.1" id="Layer_1" xmlns:x="&ns_extend;" xmlns:i="&ns_ai;" xmlns:graph="&ns_graphs;"
	 xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" width="700px" height="1000px"
	 viewBox="0 0 700 1000" enable-background="new 0 0 700 1000" xml:space="preserve">
<switch>
'''

svg_footer = '''</switch>
</svg>
'''

def getNodeTypeStr(node):
    if node.nodeType == node.ELEMENT_NODE:
        return 'ELEMENT_NODE'
    if node.nodeType == node.TEXT_NODE:
      return 'TEXT_NODE'
    if node.nodeType == node.CDATA_SECTION_NODE:
      return 'CDATA_SECTION_NODE'
    if node.nodeType == node.ENTITY_REFERENCE_NODE:
      return 'ENTITY_REFERENCE_NODE'
    if node.nodeType == node.PROCESSING_INSTRUCTION_NODE:
      return 'PROCESSING_INSTRUCTION_NODE'
    if node.nodeType == node.COMMENT_NODE:
      return 'COMMENT_NODE'
    if node.nodeType == node.NOTATION_NOD:
        return 'NOTATION_NOD'
    return None

def getChildNodeWithName(node, name):
    for n in node.childNodes:
        if n.nodeName == name:
            return n;

def createSvgFromNodes(nodes, filename):
    f = open(filename + '.svg', 'w')
    f.write(svg_header)
    if len(nodes) > 1:
        f.write('<g>')

    for n in nodes:
        if n:
            f.write(n.toxml() + '\n')

    if len(nodes) > 1:
        f.write('</g>')

    f.write(svg_footer)
    f.close()

def extractSvgContentAsNode(filename):
    doc = minidom.parse(filename)
    svgNodes = doc.getElementsByTagName('svg')
    if len(svgNodes) > 1:
        print 'Error: Unable to parse more than one svg tag from file', filename
        sys.exit(-1)

    #svgNode = svgNodes[0]
    switchNode = doc.getElementsByTagName('switch')[0]
    if switchNode.hasChildNodes:
        gnode = getChildNodeWithName(switchNode, 'g')
        if gnode and gnode.nodeType == gnode.ELEMENT_NODE and gnode.hasAttribute('id'):
            att = gnode.getAttribute('id')
            return gnode, att

        nodes = []
        atts = []
        for node in gnode.childNodes:
            #print 'child name:', node.nodeName, 'type:', getNodeTypeStr(node)
            if node.nodeType == node.ELEMENT_NODE and node.hasAttribute('id'):
                atts += [node.getAttribute('id')]
                nodes += [node]
        return nodes, atts

usage_string='''

extract player_001.svg:
    ./svgtool.py -i player_001.svg

merge extracted svgs:
    ./svgtool.py  -i player_001_extracted/Cuello.svg  player_001_extracted/Cabeza.svg player_001_extracted/Cejas.svg player_001_extracted/DientesVampiro.svg -o output
'''

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='script to extract elements from svg or merge svgs',
            formatter_class=argparse.ArgumentDefaultsHelpFormatter,
            usage=usage_string
            )
    parser.add_argument('-o', '--output', help='name of resulting svg file', required=False)
    parser.add_argument('-i', '--input', nargs='+', help='list of svg files to merge (separated by spaces). If only one file is given then, we assume that you want to extract the svg contents', required=True)

    args = parser.parse_args()

    files = args.input
    outputname = args.output

    extract = len(files) == 1

    if extract:
        inputfilename = files[0]
        if outputname:
            dirname = outputname
        else:
            dirname = os.path.splitext(inputfilename)[0] + '_extracted'
        if not os.path.exists(dirname):
            os.makedirs(dirname)
        nodes, atts = extractSvgContentAsNode(inputfilename)
        for node,att in zip(nodes,atts):
            createSvgFromNodes([node], os.path.join(dirname, att))

    else: # merge svgs
        nodes = []
        for f in files:
            node, _ = extractSvgContentAsNode(f)
            nodes += [node]
        createSvgFromNodes(nodes, outputname)
